/*
 * Copyright (c) 2023 Typist Tech Limited
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/typisttech/wpsecadvi/composer"
	"github.com/typisttech/wpsecadvi/wordfence"
	"net/http"
	"os"
)

var (
	scanner    bool
	production bool

	baseContent []byte

	wordfenceGenCmd = &cobra.Command{
		Use: "gen {--scanner|--production}",
		Example: `
  
  $ wpsecadvi wordfence gen --scanner

	Generate from Wordfence scanner feed 

  $ wpsecadvi wordfence gen --production

	Generate from Wordfence production feed

  $ wpsecadvi wordfence gen --base /path/to/composer.base.json 

	Merge with a base JSON file

  $ wpsecadvi wordfence gen --ignore UUID1 --ignore CVE-2099-0001

	Skip vulnerabilities by CVEs (production feed only) or UUIDs

  $ wpsecadvi wordfence gen --plugin-vendor foo --plugin-vendor wpackagist-plugin

	Use custom plugin vendor names.

  $ wpsecadvi wordfence gen --theme-vendor foo --theme-vendor wpackagist-theme

	Use custom theme vendor names.
`,
		Short: "Generate composer conflicts from vulnerability data feed.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !scanner && !production {
				return fmt.Errorf("error: missing feed selection. Excatly one of %s flags required", []string{"scanner", "production"})
			}

			if err := readBaseContent(); err != nil {
				return err
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			g := wordfence.NewGenerator(
				newClient(),
				newSearcher(),
			)

			j, err := g.Generate(
				viper.GetStringSlice("wordfence.gen.ignore"),
			)
			cobra.CheckErr(err)

			jBytes, err := j.Merge(baseContent)
			cobra.CheckErr(err)

			fmt.Fprintln(os.Stdout, string(jBytes))
		},
	}
)

func init() {
	wordfenceCmd.AddCommand(wordfenceGenCmd)

	// Feed selection.
	wordfenceGenCmd.Flags().BoolVar(&scanner, "scanner", false, "Generate from Wordfence scanner feed")
	wordfenceGenCmd.Flags().BoolVar(&production, "production", false, "Generate from Wordfence production feed")
	wordfenceGenCmd.MarkFlagsMutuallyExclusive("scanner", "production")

	// Ignore UUIDs and CVEs.
	defaultIgnore := []string{
		// As of 03rd January 2023, this vulnerability affects all WordPress versions.
		// https://www.wordfence.com/threat-intel/vulnerabilities/wordpress-core/wordpress-core-611-unauthenticated-blind-server-side-request-forgery
		"CVE-2022-3590", "112ed4f2-fe91-4d83-a3f7-eaf889870af4",
	}
	wordfenceGenCmd.Flags().StringArrayP("ignore", "i", defaultIgnore, "CVEs or UUIDs to ignore")
	viper.BindPFlag("wordfence.gen.ignore", wordfenceGenCmd.Flags().Lookup("ignore"))
	viper.SetDefault("wordfence.gen.ignore", defaultIgnore)

	wordfenceGenCmd.Flags().StringP("base", "b", "", "Base composer.json to merge")
	viper.BindPFlag("wordfence.gen.base", wordfenceGenCmd.Flags().Lookup("base"))

	pv := []string{composer.WPackagistPluginVendor}
	wordfenceGenCmd.Flags().StringArrayP("plugin-vendor", "p", pv, "Plugin vendor")
	viper.BindPFlag("wordfence.gen.plugin-vendor", wordfenceGenCmd.Flags().Lookup("plugin-vendor"))
	viper.SetDefault("wordfence.gen.plugin-vendor", pv)

	tv := []string{composer.WPackagistThemeVendor}
	wordfenceGenCmd.Flags().StringArrayP("theme-vendor", "t", tv, "Theme vendor")
	viper.BindPFlag("wordfence.gen.theme-vendor", wordfenceGenCmd.Flags().Lookup("theme-vendor"))
	viper.SetDefault("wordfence.gen.theme-vendor", tv)
}

func readBaseContent() error {
	base := viper.GetString("wordfence.gen.base")
	if base == "" {
		baseContent = []byte(`{}`)
		return nil
	}

	bc, err := os.ReadFile(base)
	if err != nil {
		return err
	}

	baseContent = bc

	return nil
}

func newSearcher() composer.CompositedSearcher {
	s := composer.NewCompositedSearcher()

	for _, v := range viper.GetStringSlice("wordfence.gen.plugin-vendor") {
		if v == "" {
			continue
		}
		s.AddSearcher(
			composer.NewPrefixedSearcher(composer.WPPlugin, v),
		)
	}
	for _, v := range viper.GetStringSlice("wordfence.gen.theme-vendor") {
		if v == "" {
			continue
		}
		s.AddSearcher(
			composer.NewPrefixedSearcher(composer.WPTheme, v),
		)
	}

	return s
}

func newClient() wordfence.Client {
	hc := http.DefaultClient

	if scanner {
		return wordfence.NewScannerFeedClient(hc)
	}

	return wordfence.NewProductionFeedClient(hc)
}

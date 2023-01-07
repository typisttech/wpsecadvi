/*
 * Copyright (c) 2022-2023 Typist Tech Limited
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

package composer

import (
	"bytes"
	"encoding/json"
)

func jsonUnescapedMarshal(v any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(v)

	return buf.Bytes(), err
}

// jsonMerge merges the two JSON-marshaled []byte
// preferring a over b when the keys from both objects
// are included and their values merged recursively.
//
// It returns an error if a or b cannot be JSON-unmarshalled.
func jsonMerge(a, b []byte) ([]byte, error) {
	var j1 interface{}
	err := json.Unmarshal(a, &j1)
	if err != nil {
		return nil, err
	}

	var j2 interface{}
	err = json.Unmarshal(b, &j2)
	if err != nil {
		return nil, err
	}

	merged := mergeInterfaces(j1, j2)

	return jsonUnescapedMarshal(merged)
}

func mergeInterfaces(x1, x2 interface{}) interface{} {
	switch x1 := x1.(type) {
	case map[string]interface{}:
		x2, ok := x2.(map[string]interface{})
		if !ok {
			return x1
		}
		for k, v2 := range x2 {
			if v1, ok := x1[k]; ok {
				x1[k] = mergeInterfaces(v1, v2)
			} else {
				x1[k] = v2
			}
		}
	case nil:
		// mergeInterfaces(nil, map[string]interface{...}) -> map[string]interface{...}
		x2, ok := x2.(map[string]interface{})
		if ok {
			return x2
		}
	}
	return x1
}

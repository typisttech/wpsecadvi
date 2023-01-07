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

type JSON struct {
	conflicts map[string]Link
}

type stringer interface {
	String() string
}

type Link struct {
	name        string
	constraints stringer
}

func NewLink(name string, constraints stringer) Link {
	return Link{
		name:        name,
		constraints: constraints,
	}
}

func (j *JSON) AddConflict(l Link) {
	if j.conflicts == nil {
		j.conflicts = make(map[string]Link)
	}

	j.conflicts[l.name] = l
}

type marshallableJSON struct {
	Conflict map[string]string `json:"conflicts"`
}

func (j JSON) MarshalJSON() ([]byte, error) {
	conflict := make(map[string]string, len(j.conflicts))
	for _, l := range j.conflicts {
		conflict[l.name] = l.constraints.String()
	}

	mj := marshallableJSON{
		Conflict: conflict,
	}

	return jsonUnescapedMarshal(mj)
}

func (j JSON) Merge(other []byte) ([]byte, error) {
	bs, err := j.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return jsonMerge(bs, other)
}

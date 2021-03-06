// Copyright © 2019 Máté Birkás <gadfly16@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package teflon

import (
	"errors"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// MetaSelector Nodes

// AllMetaNode returns all metadata of an object
type AllMetaNode struct{}

// NumberNode represents a number literal
type NumberNode struct {
	Value float64
}

// StringNode represents a string literal
type StringNode struct {
	Value string
}

// MetaNode represents a metadata identifier
type MetaNode struct {
	NameList []string
}

// Adds numbers numberically and concatenate strings
type AddNode struct {
	first  ENode
	second ENode
}

// Adds numbers numberically and concatenate strings
type SubNode struct {
	first  ENode
	second ENode
}

// Adds numbers numberically and concatenate strings
type MulNode struct {
	first  ENode
	second ENode
}

// Adds numbers numberically and concatenate strings
type DivNode struct {
	first  ENode
	second ENode
}

//
// ObjectSelector nodes
//

type RelPath struct {
	next   *ONode
	noMore bool
	count  int
}

type AbsPath struct {
	next   *ONode
	noMore bool
	count  int
}

type ExactName struct {
	next   *ONode
	noMore bool
	name   string
}

type MultiName struct {
	next      *ONode
	pattern   *regexp.Regexp
	index     int
	moreOfObj *TeflonObject
	children  []string
}

//
// ENode Implementations
//

func (amn *AllMetaNode) Eval(c *Context) (interface{}, error) {
	return c.IMap, nil
}

func (N *NumberNode) Eval(c *Context) (interface{}, error) {
	return N.Value, nil
}

func (S *StringNode) Eval(c *Context) (interface{}, error) {
	return S.Value, nil
}

func (M *MetaNode) Eval(c *Context) (interface{}, error) {
	var val interface{}
	v := c.IMap
	for i, n := range M.NameList {
		// Create lower map for case insensitive matching
		lm := map[string]string{}
		for k := range v {
			lm[strings.ToLower(k)] = k
		}

		var ok bool
		val, ok = v[lm[strings.ToLower(n)]]
		if !ok {
			return nil, errors.New("Couldn't find key in meta: " + n)
		}

		// If there is more name to come
		if i < len(M.NameList)-1 {
			switch val.(type) {
			case map[string]interface{}:
				// Convert next level to map
				v = val.(map[string]interface{})
			default:
				return nil, errors.New("Couldn't find key in meta: " + n)
			}
		}
	}
	return val, nil
}

func (a *AddNode) Eval(c *Context) (interface{}, error) {
	fi, err := a.first.Eval(c)
	if err != nil {
		return nil, err
	}

	si, err := a.second.Eval(c)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch f := fi.(type) {
	case float64:
		switch s := si.(type) {
		case float64:
			v = f + s
		}
	}
	return v, nil
}

func (a *SubNode) Eval(c *Context) (interface{}, error) {
	fi, err := a.first.Eval(c)
	if err != nil {
		return nil, err
	}

	si, err := a.second.Eval(c)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch f := fi.(type) {
	case float64:
		switch s := si.(type) {
		case float64:
			v = f - s
		}
	}
	return v, nil
}

func (a *MulNode) Eval(c *Context) (interface{}, error) {
	fi, err := a.first.Eval(c)
	if err != nil {
		return nil, err
	}

	si, err := a.second.Eval(c)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch f := fi.(type) {
	case float64:
		switch s := si.(type) {
		case float64:
			v = f * s
		}
	}
	return v, nil
}

func (a *DivNode) Eval(c *Context) (interface{}, error) {
	fi, err := a.first.Eval(c)
	if err != nil {
		return nil, err
	}

	si, err := a.second.Eval(c)
	if err != nil {
		return nil, err
	}

	var v interface{}

	switch f := fi.(type) {
	case float64:
		switch s := si.(type) {
		case float64:
			v = f / s
		}
	}
	return v, nil
}

//
// ONode Implementations
//

// RelPath

func (rpn *RelPath) NextMatch(o *TeflonObject) (res *TeflonObject) {
	var err error
	if rpn.noMore {
		rpn.noMore = false
		return nil
	}

	// Give back o or traverse upvards.
	if rpn.count > 1 {
		fspath := o.Path
		for i := 1; i < rpn.count; i++ {
			fspath = filepath.Dir(fspath)
		}
		res, err = NewTeflonObject(fspath)
		if err != nil {
			log.Fatalln("FATAL: Couldn't create object:", fspath)
		}
	} else {
		res = o
	}

	if rpn.next == nil {
		rpn.noMore = true
	} else {
		res = (*rpn.next).NextMatch(res)
		if res == nil {
			rpn.noMore = true
		}
	}

	return res
}

func (node *RelPath) GenerateAll(fspSl []string) (res []string) {
	for _, fsp := range fspSl {
		// Give back o or traverse upvards.
		for i := 1; i < node.count; i++ {
			fsp = filepath.Dir(fsp)
		}
		res = append(res, fsp)
	}

	if node.next != nil {
		res = (*node.next).GenerateAll(res)
	}

	return res
}

func (rpn *RelPath) SetNext(node *ONode) {
	rpn.next = node
}

// AbsPath

func (apn *AbsPath) NextMatch(o *TeflonObject) (res *TeflonObject) {
	var err error
	if apn.noMore {
		return nil
	}

	if apn.count == 1 {
		res, err = NewTeflonObject("/")
		if err != nil {
			return nil
		}
	} else {
		res, err = NewTeflonObject("//")
		if err != nil {
			return nil
		}
	}

	if apn.next == nil {
		// In itself AbsPath reurns only one match.
		apn.noMore = true
	} else {
		res = (*apn.next).NextMatch(res)
		if res == nil {
			apn.noMore = true
		}
	}

	return res
}

func (node *AbsPath) GenerateAll(fspSl []string) (res []string) {
	if node.count == 1 {
		res = []string{"/"}
	} else {
		shw, err := NewTeflonObject("//")
		if err != nil {
			return []string{}
		}
		res = []string{shw.Path}
	}

	if node.next != nil {
		res = (*node.next).GenerateAll(res)
	}
	return res
}

func (apn *AbsPath) SetNext(node *ONode) {
	apn.next = node
}

// ExactName

func (enn *ExactName) NextMatch(o *TeflonObject) (res *TeflonObject) {
	var err error
	if enn.noMore {
		return nil
	}

	res, err = NewTeflonObject(filepath.Join(o.Path, enn.name))
	if err != nil {
		enn.noMore = true
		return nil
	}

	if enn.next == nil {
		enn.noMore = true
	} else {
		res = (*enn.next).NextMatch(res)
		if res == nil {
			enn.noMore = true
		}
	}

	return res
}

func (node *ExactName) GenerateAll(fspSl []string) (res []string) {
	for _, fsp := range fspSl {
		fsp = filepath.Join(fsp, node.name)
		res = append(res, fsp)
	}

	if node.next != nil {
		res = (*node.next).GenerateAll(res)
	}

	return res
}

func (enn *ExactName) SetNext(node *ONode) {
	enn.next = node
}

// MultiName

func (mnn *MultiName) NextMatch(o *TeflonObject) (res *TeflonObject) {
	var err error

	// Init children.
	if mnn.children == nil {
		mnn.children = o.ChildrenNames()
	}

	for res == nil {
		if mnn.moreOfObj != nil {
			res = mnn.moreOfObj
		} else {
			log.Printf("DEBUG: Looking for multi matches: mnn.index: %v", mnn.index)
			for mnn.index < len(mnn.children) {
				name := mnn.children[mnn.index]
				mnn.index++
				if mnn.pattern.MatchString(name) {
					log.Printf("DEBUG: Found match: name: %v", name)
					res, err = NewTeflonObject(filepath.Join(o.Path, name))
					if err != nil {
						log.Fatalln("FATAL: Couldn't create found object.", name)
					}
					break
				}
			}

			// No more matching child.
			if res == nil {
				mnn.index = 0
				mnn.moreOfObj = nil
				mnn.children = nil
				return nil
			}
		}

		if mnn.next != nil {
			log.Printf("DEBUG: Calling child of %v", res.FileInfo.Name)
			mnn.moreOfObj = res
			res = (*mnn.next).NextMatch(res)
			if res == nil {
				mnn.moreOfObj = nil
			}
		}
	}
	return res
}

func (node *MultiName) GenerateAll(fspSl []string) (res []string) {
	for _, fsp := range fspSl {
		o, err := NewTeflonObject(fsp)
		if err != nil {
			continue
		}
		for _, ch := range o.ChildrenNames() {
			if node.pattern.MatchString(ch) {
				res = append(res, filepath.Join(fsp, ch))
			}
		}
	}

	if node.next != nil {
		res = (*node.next).GenerateAll(res)
	}

	return res
}

func (mnn *MultiName) SetNext(node *ONode) {
	mnn.next = node
}

//
// Utility Functions
//

// NumberNode needs to be Stringer for string concatenation
func (N NumberNode) String() string {
	return strconv.FormatFloat(N.Value, 'G', -1, 64)
}

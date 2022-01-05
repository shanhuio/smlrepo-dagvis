// Copyright (C) 2021  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package dagvis provides clear visualization of a DAG.
package dagvis

import (
	"flag"
	"log"

	"shanhu.io/aries"
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/osutil"
)

type server struct {
	static *aries.StaticFiles
	tmpls  *aries.Templates
}

func (s *server) serveIndex(c *aries.C) error {
	dat := struct{}{}
	return s.tmpls.Serve(c, "dagview.html", &dat)
}

func makeService(home string) (aries.Service, error) {
	h, err := osutil.NewHome(home)
	if err != nil {
		return nil, errcode.Annotate(err, "make new home")
	}

	s := &server{
		static: aries.NewStaticFiles(h.Lib("static")),
		tmpls:  aries.NewTemplates(h.Lib("tmpl"), nil),
	}

	serveStatic := s.static.Serve

	r := aries.NewRouter()
	r.Index(s.serveIndex)
	r.Get("style.css", serveStatic)
	r.Dir("js", serveStatic)
	r.Dir("jslib", serveStatic)

	return r, nil
}

// Main is main.
func Main() {
	addr := aries.DeclareAddrFlag("")
	home := flag.String("home", ".", "home dir")
	flag.Parse()

	s, err := makeService(*home)
	if err != nil {
		log.Fatal(err)
	}
	if err := aries.ListenAndServe(*addr, s); err != nil {
		log.Fatal(err)
	}
}

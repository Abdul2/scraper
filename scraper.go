// Copyright 2016 abdulrashid2@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main
//this code is designed to scrape data from a defined target website
//it can not be re-used without modifications
//key to any data scrapping is the full understanding how
// to navigate through the content and reach data relevant to your case
// at all times, it is advised to check the suitability of collecting
//information from website that does not belong to you.
//check with the right people first.
import (

	"fmt"
	"net/http"
	"golang.org/x/net/html"

	"regexp"
	"strings"
	"runtime/debug"
	"encoding/json"

)

//countlines is used to number raws to make it easier to reference and work
//on collecting information
var countlines int = 0

//Raw struct will hold line number, content and the game date line number.
//Raw will have to be changed depending on what information is collected.
type Raw struct {

	Linenumber int
	Content string
	Gamedatelinenumber int
}

var raws = make([]Raw,0)
var gamedatelocations = make([]int,0)  // to hold line numbers where game dates are shown


type Fixture struct {

	Gamedate string `json:"gamedate"`
	Team string  `json:"team"`
	Awayorhome string `json:"awayorhome"`
	Competition string `json:"competition"`
	Result string `json:"result"`
	Score1 string `json:"score1"`
	Score2 string `json:"score2"`
}

var season []Fixture

//main programme entry
//the programme uses net.http to make a Get call and tokenises the returned body
// using html tokenizer to reach the content
// the content is manipulated to produce the desired structure
func main() {

	//regexp is used to define data of interest
	gamedatepattern, err := regexp.Compile("[aA-zZ]{3}\\s[0-9]{2}\\s[aA-zZ]{3}\\s[0-9]{4}") //Sat 22 Aug 2015
	if err != nil {
		f(err)
	}

	//test locally first on representative gape
	resp, err1 := http.Get("http://localhost:8080/outputfile1946_1947.html")
	if err != nil {
		f(err1)
	}

	t := html.NewTokenizer(resp.Body)
	countlines = 0

	L:        //Label to be used to break out from the combined switch and for loop


	// walk through page content and produce two slices.
	//raws slice holds describes each raw and include raw number, content and 0 or 1
	// value to indicate if the raw holds data of interest
	// gamedatelocation slice holds raw numbers of where game dates are shown
	for {

		tt := t.Next()
		switch  {

		case tt == html.ErrorToken:
			break L

		case tt == html.TextToken:

			line := strings.TrimSpace(string(t.Text()))
			raw := Raw{

				Linenumber: countlines,
				Content: line,
				Gamedatelinenumber: 0,
			}

			raws = append(raws, raw)

			if gamedatepattern.MatchString(string(line)) {

				raws[countlines].Gamedatelinenumber = 1

			}

			countlines++
		}

	}



	var fixture Fixture

	for i := 0; i < len(raws); i++ {

		if raws[i].Gamedatelinenumber == 1 {

			for c := i; c < i + 9; c++ {

				rawnumber := raws[c].Linenumber
				rawcontent := raws[c].Content

				if rawcontent != "" { //ignore empty raws

					fmt.Println(rawnumber, rawcontent)

					j := i

					fixture.Gamedate = raws[j].Content
					fixture.Team = raws[j+2].Content
					fixture.Awayorhome = raws[j+4].Content
					fixture.Competition = raws[j+6].Content
					fixture.Result = raws[j+8].Content


				}

			}

			fmt.Print("\n---------------------\n")
			season = append(season,fixture)


		}

	}

	//marshal into `json:"


	j,_ := json.MarshalIndent (season,"","     ")



	fmt.Printf("%s",j)



}




func f(r error){



	fmt.Println("not good ")
	debug.PrintStack()


}

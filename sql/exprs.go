// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sql

// --------------------------------------------------
//
// --------------------------------------------------

func (p *Parser) parseTable() (one *Table, err error) {

	var lit string

	one = &Table{}

	_, lit, err = p.shouldBe(IDENT)
	if err != nil {
		return nil, &ParseError{Found: lit, Expected: []string{"table name"}}
	}
	one.Name = lit

	return one, err

}

func (p *Parser) parseTables() (mul []Expr, err error) {

	for {

		one, err := p.parseTable()
		if err != nil {
			return nil, err
		}

		mul = append(mul, one)

		// If the next token is not a comma then break the loop.
		if _, _, exi := p.mightBe(COMMA); !exi {
			break
		}

	}

	return

}

// --------------------------------------------------
//
// --------------------------------------------------

func (p *Parser) parseThing() (one *Thing, err error) {

	var tok Token

	one = &Thing{}

	_, _, err = p.shouldBe(EAT)
	if err != nil {
		return nil, &ParseError{Found: one.Table, Expected: []string{"@"}}
	}

	tok, one.Table, err = p.shouldBe(IDENT, NUMBER)
	if err != nil {
		return nil, &ParseError{Found: one.Table, Expected: []string{"table name"}}
	}

	_, _, err = p.shouldBe(COLON)
	if err != nil {
		return nil, err
	}

	tok, one.Thing, err = p.shouldBe(IDENT, NUMBER)
	if err != nil {
		return nil, &ParseError{Found: one.Thing, Expected: []string{"table id"}}
	}

	val, err := declare(tok, one.Thing)

	one.ID = val

	return

}

func (p *Parser) parseThings() (mul []Expr, err error) {

	for {

		// var one *Thing

		one, err := p.parseThing()
		if err != nil {
			return nil, err
		}

		mul = append(mul, one)

		// If the next token is not a comma then break the loop.
		if _, _, exi := p.mightBe(COMMA); !exi {
			break
		}

	}

	return

}

// --------------------------------------------------
//
// --------------------------------------------------

func (p *Parser) parseIdent() (*IdentLiteral, error) {

	tok, lit, err := p.shouldBe(IDENT)
	if err != nil {
		return nil, &ParseError{Found: lit, Expected: []string{"name"}}
	}

	val, err := declare(tok, lit)

	return val.(*IdentLiteral), err

}

func (p *Parser) parseNumber() (*NumberLiteral, error) {

	tok, lit, err := p.shouldBe(NUMBER)
	if err != nil {
		return nil, &ParseError{Found: lit, Expected: []string{"number"}}
	}

	val, err := declare(tok, lit)

	return val.(*NumberLiteral), err

}

func (p *Parser) parseString() (*StringLiteral, error) {

	tok, lit, err := p.shouldBe(STRING)
	if err != nil {
		return nil, &ParseError{Found: lit, Expected: []string{"string"}}
	}

	val, err := declare(tok, lit)

	return val.(*StringLiteral), err

}

func (p *Parser) parseRegion() (*StringLiteral, error) {

	tok, lit, err := p.shouldBe(IDENT, STRING, REGION)
	if err != nil {
		return nil, &ParseError{Found: lit, Expected: []string{"string"}}
	}

	val, err := declare(tok, lit)

	return val.(*StringLiteral), err

}

func (p *Parser) parseBoolean() (*BooleanLiteral, error) {

	tok, lit, err := p.shouldBe(TRUE, FALSE)
	if err != nil {
		return nil, &ParseError{Found: lit, Expected: []string{"boolean"}}
	}

	val, err := declare(tok, lit)

	return val.(*BooleanLiteral), err

}

func (p *Parser) parseDefault() (Expr, error) {

	tok, lit, err := p.shouldBe(TRUE, FALSE, NUMBER, STRING, REGION, ARRAY, JSON)
	if err != nil {
		return nil, err
	}

	return declare(tok, lit)

}

func (p *Parser) parseExpr() (mul []*Field, err error) {

	var tok Token
	var lit string
	var exi bool

	for {

		one := &Field{}

		tok, lit, err = p.shouldBe(IDENT, ID, NULL, ALL, TIME, TRUE, FALSE, STRING, REGION, NUMBER, DOUBLE, JSON)
		if err != nil {
			return nil, &ParseError{Found: lit, Expected: []string{"field name"}}
		}

		one.Expr, err = declare(tok, lit)
		if err != nil {
			return
		}

		// Next token might be AS
		if _, _, exi = p.mightBe(AS); exi {

			tok, lit, err = p.shouldBe(IDENT)
			if err != nil {
				return nil, &ParseError{Found: lit, Expected: []string{"field alias"}}
			}

			one.Alias, err = declare(tok, lit)
			if err != nil {
				return
			}

		}

		mul = append(mul, one)

		// If the next token is not a comma then break the loop.
		if _, _, exi = p.mightBe(COMMA); !exi {
			break
		}

	}

	return

}

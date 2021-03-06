package main

import (
	"fmt"
	"sort"
	"strings"
)

var openers = map[string]bool{
	"(": true,
	"{": true,
	"[": true,
	"<": true,
}

var closers = map[string]int{
	")": 3,
	"]": 57,
	"}": 1197,
	">": 25137,
}

type delim string

const (
	parens delim = "()"
	square delim = "[]"
	curly  delim = "{}"
	angle  delim = "<>"
)

type stackEntry struct {
	opener delim
}

func main() {
	incompleteStacks := [][]stackEntry{}

	for _, l := range strings.Split(input, "\n") {
		stack := []stackEntry{}
		isInvalid := false
		var invalidCloser delim
	charLoop:
		for _, c := range strings.Split(l, "") {
			validCloser := ""
			if len(stack) > 0 {
				// fmt.Printf("wat?> %+v\n", stack[len(stack)-1].opener)
				switch stack[len(stack)-1].opener {
				case parens:
					validCloser = ")"
				case square:
					validCloser = "]"
				case curly:
					validCloser = "}"
				case angle:
					validCloser = ">"
				}
			}
			// fmt.Printf("valid char: %q\n", validCloser)
			switch c {
			case "(":
				stack = append(stack, stackEntry{opener: parens})
			case "[":
				stack = append(stack, stackEntry{opener: square})
			case "{":
				stack = append(stack, stackEntry{opener: curly})
			case "<":
				stack = append(stack, stackEntry{opener: angle})
			case validCloser:
				// stack = append(stack, stackEntry{closer: delim(validCloser)})
				stack = stack[0 : len(stack)-1]
			default:
				// fmt.Printf("Invalid char found\n")
				isInvalid = true
				invalidCloser = delim(c)
				break charLoop
			}
			// fmt.Printf("Stack: %+v\n", stack)
		}
		if !isInvalid {
			incompleteStacks = append(incompleteStacks, stack)
		}
		fmt.Printf("Line %s: %t, %v\n", l, !isInvalid, invalidCloser)
	}

	sums := []int{}
	fmt.Printf("Incompletes: %+v\n", len(incompleteStacks))
	for _, stack := range incompleteStacks {
		subSum := 0
		for i := len(stack) - 1; i >= 0; i-- {
			s := stack[i].opener
			switch s {
			case parens:
				subSum = subSum*5 + 1
			case square:
				subSum = subSum*5 + 2
			case curly:
				subSum = subSum*5 + 3
			case angle:
				subSum = subSum*5 + 4
			default:
				fmt.Printf("Shouldn't happen: %q\n", s)
			}
		}
		sums = append(sums, subSum)
		fmt.Printf("subSum: %d\n", subSum)
	}
	sort.Ints(sums)
	mid := len(sums) / 2
	fmt.Printf("Total: %d\n", sums[mid])
}

var sample = `[({(<(())[]>[[{[]{<()<>>
[(()[<>])]({[<{<<[]>>(
{([(<{}[<>[]}>{[]{[(<()>
(((({<>}<{<{<>}{[]{[]{}
[[<[([]))<([[{}[[()]]]
[{[{({}]{}}([{[{{{}}([]
{<[[]]>}<{[{[{[]{()[[[]
[<(<(<(<{}))><([]([]()
<{([([[(<>()){}]>(<<{{
<{([{{}}[<[[[<>{}]]]>[]]`

var input = `({[[{[(([({((<[]>{[][]})[{<>[]}{[]{}}])[<[<><>][[]{}]>{[{}()][{}]}]}(<[{<>{}}[()]]<[{}[]]>><[{()}({}())]
[[<([<<<((([<<()()>{{}()}>[[()<>]]]{<{<>[]}(()<>)>[<{}{}>[{}]]}){{({<>()}{[]})<[[][]]>}[{[
[{{<{({<<({{{<[]><<>[]>}{{()()}(<>[])}}<[{(){}}<<>{}>]((<>[]){{}{}})>}([[{[]<>]][(<>[])(()<>)]])){[
[<([({{{<{[[[{{}{}}]{<[]()>([]}}][({()()}){<()<>><(){}>}]]<<(<[][]>({}{}))([[]<>]<[]()>)>{({<>[]}{
[{(({<({(<(({<[]{}>(<>())}[([]()){()<>}]){<<(){}>{()<>}>[{<>{}}{<>()}]})({[<<>()><{}[]>]})>)[{([[{(){
(<(({[<{[{[{<[()<>](()[])>}[<[()<>]((){})>([<>[]][[]{}])]]([[[[][]](()[])]<(<><>)<[]<>>>][{(<>{})}<<[]{}>[{}(
<[[<[{[<[[<<([{}[]][{}<>])<<<>()>([]{})>><((())[(){}])({<>()}{[]()})}>({(([]{})<{}[]>){{(){}}[[][]]}}
<{(<<<<{{({{((()())<<>{}>)([<>()][[]()])}[[{[][]}]<<()<>>[()<>]>]})}<<[<<({}<>>[[]<>]><[()[]]{{}{}}>>]<<
<[[[{<({(<{[((()()))]<{<()>{<><>}}((<>{})[{}[]])>}>)>)<((<<[{([]<>)([]<>)}[[[]{}](()[])]]({{()[]}})>{
{[{[<<[<[<{[{{{}[]}[[]()]}{[[]<>]([]<>)}]<((()<>){<>()}){<[]()>{[]()}}>})<[{<{[]{}}[<>]><([]
<(<{[[{((<({((<>())<()<>>)<(<><>)<(){}>>}]({<[{}{}]{()<>}><<()<>>(<>())>}[{[<>[]][<><>]}<({}()){()()}>
[(([<{<[({{[<[(){}](()[])>]<<<[]()>{{}[]}><{()}(<>())>>}[<[{{}()}[<>{}]]>(<[<>[]]<()[]>>{<()><{}<>
{<[<<[{[([[[{<{}<>>{()[]}}[<[]()>([]<>)]][<{[]{}}{[]<>}>]]{[([(){}]{()<>}>{[<>{}][<>()]}]{<[{}
[[<{<<{{{(<({[<>[]]}{<(){}>[<>[]]}){([[]{}])({[]{}}(<>[]))}><<<{[]()}{<>()}>{(<>)<()[]>}>[[({}<>)((
<(<(({<<[(<[{([]()){[][]}}<(()[])>]><[[{()<>}[()[]]]({{}<>}{<>()})][(<[]()>[{}()])[([]{})[
(<[{<<[{<[([([()[]][[]{}])]{{(()<>)}[[(){}](()<>)]})]>([(([([][])[<>[]]]<<[]<>>>)[{<()[]>(()[])}({<
({(([<<<<<((<{{}{}}<()<>>>[<[][]>[{}<>]])([(<>()){()[]}]<<<><>>[{}<>]>>){([[{}{}][[]()]]<<[][]><[]
(((<{({{([[<<<(){}><{}<>>>>[{{[]<>}(()[])}(<[][]>[()])]]]<(({{{}[]}{<>{}}}[{<>{}}[()()]])[<[{}<>](()
<(([<<{[[{{{{(()[])([]{})}}{<<{}{}>{{}<>}>[{{}{}}<()()>]}}[{[((){})((){}]]}{<(<>[])<{}[]>>[<[]<>>({}<>)]
{([{([{[{<([<<{}()>[(){}]>(<<>{}>[<>{}])])[({[[]()]<[]{}>}<{{}{}}>)<[{<><>}](<[]<>>[{}<>])>]>}]}])[{{{[
[(<{(([{{{({([[][]][()[]])}<<<()[]><()>><(<>{})<<><>>>>)<{[[<><>]<{}[]>]}>}([(<{<><>}{{}[]}>[[<>()](<><>)
(<(((<{(<[({([<>[]]{{}{}}){[[]{}][()<>]}}{[{(){}}[()<>]]<<<>[]>([][])>})][<{<[[]<>]((){})>([<>()]<(){}>)}>]>)
[{{[(({(<<[({{{}{}}<[][]>}([<>][[]{}]))]<<((<>())<()[]>)<({}())<{}[]>>>>>>)}{<<[{[[(<>()){<><>}]<(
<[[[((<{[<{(<[[]()]>){[{()<>}[<>()]]<{<><>}[[]<>]>}}[<([{}{}])({[]<>}[[]])><({()[]}(<>{}))[<
<{<<[<[<(([<{<(){}>({}())}{<<>>[[][]]}>[(<{}()>{(){}})(<[]><{}{}>)]]<((((){}){<><>})<{{}[]}
[(<[<(<{({{<(<{}{}>{{}()})<{[][]}{{}[]}>>(<((){})(()[])>{(<>())})}})<<<[[([]<>)([]<>)]{(<>
(({<[<{{<{<[(({}())[[]{}])<{<>[]}(()<>)>]([(<>[])<[]()>]([[]()){{}()}))>}>[({<(([][])({}{}))<{{}[]}[<>
((({({{[[(<{<<[]()><[]{}>><(()[])[{}()]>}[{[<>()]([])}{{{}{}}<()()]}]>)<{<[[<>{}][[]()]]({<>()}{<><>})>
{[[[[((<<<{[{([][]){<>{}}}[<<>()>([]{})]](<((){})<{}()>><([]())[<><>}>)}{{[[(){}]]{<<>[]><[]{}>}}<({(){
({[[[(({<({[{[<>[]][{}{}]}<<<>>{<>{}}>]<{[{}<>]{()[]}}>}<<([[]{}]({}{}])>{[[<>[]]({}<>)]([{}[]]
{{{({{[{({{<{{()[]}{{}[]}}{<{}{}>}>{(({}())[()<>])((<>{})([]))}}}]<<[(((()<>)<(){}>){[[][]]{{}
{({<{<[{({[{{<<>>[{}{}]}[[<>[]]<{}[]>]}[<[()()]{[]}><([]<>){<>[]}>]]})}[({<[<<{}<>>(()[])>]<({{}(
[{((<{<({[{<[(()()){<><>}](<{}{}>{[][]})>[{<{}()><[]()>}]}][(<({[]{}})>)<<[({}{}){[]()}]{((){}){()[
{<<<{({[(<{<{<()<>>([][]]}<(<>[])[()[]]>>}>)]}[{[{[{[<{}[]>]({()[]}<()()>)}][[[[{}()]{<>{}}][[[
{[{{{{[((<{[[([]<>){{}[]}][<[]{}>({}{})]][[{<>()}[()[]]]<<[]{}><<>{}>>]}{{{[[]<>]<{}()>}[[[]()]{<>()
(([<[[[{{((<(<(){}>[[]<>]){{<><>}([][])}>{[{{}[]}{()[]}]<<{}<>><[]>>}))}({({[[<><>]<[]{}>][{<
[(<({[<[<<(<{([]{})(<>())}<{{}[]}{[]{}}>><(([]<>)[(){}])>)[<{[[]{}]([]{})}[<[][]>(()<>)]>{<{<>()}(()<>)>(<(
[{<(<((<{[(<{<<>()>}<<<>[]><<><>>>><{<()[]><()()>}{[(){}]{[]{}}}>)<{([<>{}][()>)<{[]{}}([][])>}>]({{<({}
({(<[<<{<(((<[{}<>]<<>()>><([]<>)>)[[<<>>{{}<>}](<{}<>>)])[{<[{}]<{}{}>><<[]<>>{[][]}>}}){(({(
(([([{{<<<{({[<>[]]<<>>}({()<>}))({(<>{})(()[])))}((<([]<>){<>[]}>(<[][]>{()<>})){([()()]{(){}})((<>())<{
[{{[<<<((<[([{()<>}<()[]>][(()<>)])({<()<>>({}<>)}{{()<>}{()()}})]>[(<([()[]]<(){}>){[{}<>]{{}{}}}>[([
[(({<[[[<<[[{[()()][{}()]}[[()<>]([]<>)]][<<<><>>[[]<>]><[()<>]<<>()>>]]([{{()[]}{()()}}]{([()[]]({}))<<[]
<[<(<(({<([<<[[]()]<[]<>>>[{[]()}<(){}>]>{<[{}[]]{[]{}}>[<{}<>>]}][<(<{}<>>{()<>})[{<><>}{[]()}]
[{[<<{{{{{{{[{[]<>}{[]{}}]}[[(<>())([][])]<({}[]}[{}[]]>]}{{<<[]<>>[[]<>]>}[[([]{})({}{})]([
({[({[[{{[{[[{<>[]}(()<>)][<()()><[][]>]]{{{<>[]}<<>[]>}<[<>{}]<[][]>>}}<(<{[]{}}[{}{})>([{}{}](<>{})))([({}
<[(<{({<([<[{<(){}><<>[]>}{[[][]]<[][]>}]{<([]{}){[]<>}>}><[[<<>{}><{}[]>]]{<<()<>>{<>()}>({{}<>}<<>[]>)}
<[({({[{[<<{<([]())({}())><[()[]](<>{})>}>{{[(<>{})({}{})]{[[][]]{()[]}}}<{[<>[]]{<>()}}[{<>()}([][]
((({([<[<{{[<<{}[])[{}{}]>({{}()}{()<>})]{<<{}()><<>()>>(((){}){[]<>})}}([([()()]{<>{}}){({}[])[{}
(<((<([{<(<<{<<>[]>[{}{}]}[{<>}([]<>)]>(<{[]{}}{{}()}><(<>[]){()<>}])>{({[()<>]{<>{}}}[[[]<>]{<><>}]){
<({<[[([<<[(<[()[]]>)({(()())(<>{})})]<(<[{}()]{{}[]}>)>>{[[{<<>()>[<>()]}{([]){<>()}}][<<[][]><(){}
((({<[<[{<(([({}()][<>()]]))[<[[()]<{}<>>]({[]()}(<>[]))>]>(<{{<<><>><()()>}({[]()}([]()))}(<{{}()}((){})>({{
[{<<{(<{{{{[{<()[]>(<>())}](({[]{}}{(){}})<({}())[[]<>]>)}{<{[[][]]>{[()()]<<>()>}>[(<{}[]>[()()]
({[((<{<<[(<[{(){}}([]<>)]<(<>{})(<>())>>{{(<><>)(<>())}[<()><{}{}>]}){{{{<>()}<()()>}{[[]
<[[[{([{((<<({[][]}<()[]>)>><[(<<>{}>((){}))]{[<[]{}>({}<>)]({()[]}(<>{}))}>)<({{([]()){{}{}}}})<({<<>{}>{(
<<({<<[[<((<<{{}()}([]()]><{()<>}<<>()>>>({<<>{}>({}<>)}{<[]<>>[<><>]})))<[[<<<><>><()()>>{<<>{}>}]([<
([<{{<(<{{(<{[[]()][{}<>]}[(<><>)]>)}{({[(<>[]){{}{}}]<{{}{}}[(){}]>})[{({(){}}(<>()>)<[[][]][{}<
<[<[([[<<[([{[()<>]<[]<>>}<{()}[()]>])]<{<[<<>[]>[[]()]]<([][]}{<>[]}>>[<<<>{}>>{<{}[]>({}())}]}<[<<[]<>>><
<{<([[(([{<<[{{}()}[<>()]}>{<[<><>](<>{})><<<>[]>>}>(([{()()}(<>())]<[<>{}]>))}]{{{[([[]<>
[{<({<<[((<{<{{}<>}[[][]]>[[{}()]<()[]>]}>{{{{[][]}<[][]>}{({}<>)[(){}]}}(([[]]{()[]})[[()<>]{
<([(<<(([{{<({<><>}{()}}{{<>()}[{}<>]}>[<({}[])(<><>)><{{}[]}((){})>]}<(<{<>{}}({}())>)>}[{(
[{<(<({[[({({<(){}>[[]{}]}({<><>}[[]{}])){[<<>[]><()>]<[{}()]({})>}}[[({{}[]}<{}[]>]]<{<<>[]>(<><>)
[({{<[[<({<{(({}<>)[()[]]){{<>{}}<()[]>}}{(<<>[]>{<>[]})({<>[]}<()[]>)}>><{(<[()][(){}]>){((<>[]){<><>}
({{[{[{<<{((({(){}}){[<>{}]}))<<<{{}[]}{<>[]}>((()[]))>{{<{}{}><()>}[{{}{}}<<><>>]}>}><[(<[((){}){<>()
<[{<([[<[[[{<[<>[]]({}<>)>[{<>[]}[{}{}]]}<{<<>{}><{}>}[{{}()}[[]{}]]>]([(((){})[()[]])[<(){}><()
({{<<<([[([[<(<>{})([][]>>]]([[(()[]){{}()}]][{<()()>[<>{}]}[{[]{}}((){})]]))(<<{<{}><[]>}(<<>()>({}<>))>[
{<{<<[<(([[[<{[]{}}[{}{}]>]{<<()[]>[(){}]>}]]<({<<()[]>{<>[]}>}{<[{}{}]{<>{}}>{[<>()}}})<<<
(({[<[[{<<<({((){})}[{()}({}[])])>>>{{{({(()<>)[[]<>]}[<[]<>>[{}[]]]){[<()[]>][(<><>)]}}}<(((<[]()
<{[({{{<([{[{{{}<>}}]<<[[]]({}{})>>}[<({[][]}<{}[]>)<{{}[]}>>]])>(((((<{()()}<<>[]>><[(){}][<>{
<<(([<<[[[[({<<><>>[{}]}[{{}[]}{[][]}])][((((){}))[<(){}>])[({()()}[[]()))]]]{{[((<>)[{}{}])({{}[]
<<({<{{(<(({<((){}){[]()}>}){{{[{}()]}{(())<[][]>}}<[[()]](<{}()>{{}{}})>}]><<[({{<>[]}[(){}
([{{<[(<[([<[{{}{}}{<>[]}][(<>[])]>](([[()<>]([]())]((()()){{}[]}))<([()<>]<<><>>){<<>{}>(<>{})}>))<{{{<()[
{[[[[([{[{[{{<<>()><(){}>}({<>()}[{}[]])}([{(){}}<{}{}>]({<>[]}(<>[])))]}({[[[()<>]{[][]}]]}<{<{<>()}
[<(([<<((<(((<[]()>{[]{}})<{()<>}([]{})>))<(([(){}]<[]{}>)<<{}>({}())>)((<{}()>{()<>}){({}{})<(){}>})>>((<<(<
[([[({<([(<((({}()))[{{}<>}{()[]}]){{{{}()}(<>{})}<<{}()>{<>{}}>}>)<([[[{}](<>{})]<[[]()]<<>[]>>]
(({(<<[<{(([((()<>)([]()))])<{(([]<>){()<>})[{{}<>}[{}<>]]}>)({([{<><>}({}())])})}(<([([<>{}])<<<
<(([([{[[{{<[<()>](({}()){<>[]})>[[({}[])]{<()<>>{{}}}]}[<{{[]<>}[[]()]}>]}]](([<<{((){})(<>)}{[()[]][(
({(({<<[{{[{<<()[]>{[]()}>([<><>]<{}>)}<(<{}()>[[]()])[<[]}[()]]>]}}{[{(([()()][{}<>]))}<[([[][]][<>{
{(([({<<[[{<(<[]<>>([]<>)][<<>()>[[]<>]]><{(())[(){}]}<(()[]){{}[]}>>}[<[(<><>)(<>[])]<{[]}>><[[()[]][[]]]((<
{{(<[[({(([[(<{}<>>(()())>{<(){}>{<>{}}}]{[{()[]}{[]{}}]}]<[[[<><>][[]<>]][{{}{}}]]<([{}{}]<<><>>)<
[<<{<(<[([[(<<()<>>([]{})>(<{}<>>{[]<>})){<[<>{}]<[][]>>(<()[]><<><>>)}]<([(<>()){<>{}}](<
([[{{([[([{[<<()<>>[()[]}>[[<>[]](<><>)]][{{{}()}{{}{}}}(<(){}>[{}()])]}]){{{{{((){})[{}<>]}<{{}<>}[[]
[[{(({{{<[(<{(<>())<[]<>>}[([][])]>)][<[({()()}(<>())){[{}<>]<[]()>}]}[[{{(){}}{()()}}]{[{{}{}}[{
({<{<{[<<{({{{{}{}}(<>())}<{[]}<{}{}>>}(([[]{}]([]{}))))}>>]}<{{<({[<([]())>(({}<>)<{}<>>)]}<[<(()()}{{}()}
<<[[<[<[[{({<[<>[]][()()]>{<{}()>{<><>}}}[{([][])<()<>>}<{{}()}([][])>])[{[[{}]({}{})]{{[][]}([][])}}
[({[[[[[[({{{<()<>>}}[[(()())({}[])][{[]{}}[()()]]]}{{{{[]()}<[]>}[[<>[]]<<><>>]}<[<{}<>>(<>)]<{<>{}
{[[([([({{[{<{{}[]}{()[]}>}{({<>()})<<<>{}><{}{})>}][[<[{}{}][()<>]>{[[][]][<>[]]}]{[{()[]}[{}()]]({[]<>
{{[[<{({({[[<<<>{}><{}[]>><([][])(()<>)>]]}<([{[<>]<<>{}>}<[{}[]]([][])>][<({}{}>([]<>)>{[<>{}]}])[
({{[{{(({{<({(<><>)<[]{}>}[[{}[]]<[]{}>])>(([{{}()}{[][]}]{{()()}<{}<>>}){<{{}[]}{{}{}}>{(()())}})
{({<{[<<<{<<<{{}()}[[][]]>{<()>(<>())}>>({[[<>[]]<(){}>]<[[]{}]{<><>}>}<[(()[])<<>[]>][<[]()>]>)}({<
{<{([[<[<{{[[[{}<>]<(){}]][<()>(<>{})]]<<<()[]>[<>[]]>>}}{(<{{[]<>}<<>[]>}(<[]<>>)>([{()()}[[]<>
<<((<<<<{<<{({<>[]}({}<>))}(([<><>]{<>()}>([<>[]][(){}]))>><(({<()[]>}({[]<>}[<>[]]))<({(){}}
[{<({<<<{({[{[[]()]<<>>}{{[]()}[<><>]}]}{(({<><>}<[][]>){({}{})(<>{})})})}{<<({(()())<[]()]})([{<><>}<()<
(<((([(<<<<<{[{}{}][{}{}]}(<<><>>[[]<>])>{<{<>{}}<<>()>><{<>[]}[{}{}]>}>(((<<><>>({}())))([<<>[]>]<<{}()>([]
<[{([<<<{{[<<({}())[<>]>[{[]{}}[()]]>{{[<>{}]{<><>}}<({}[]){()<>}>}]}}[{{({{()()}(()[])}<<[`

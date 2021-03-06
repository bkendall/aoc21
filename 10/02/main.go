package main

import (
	"fmt"
	"sort"
	"strings"
)

var closeForOpen = map[string]string{
	"(": ")",
	"{": "}",
	"[": "]",
	"<": ">",
}

var finalPoints = map[string]int{
	")": 1,
	"]": 2,
	"}": 3,
	">": 4,
}

func main() {
	lines := strings.Split(input, "\n")

	incompleteStacks := [][]string{}

	for _, l := range lines {
		stack := []string{}
		isValid := true
		for _, c := range strings.Split(l, "") {
			if !isValid {
				continue
			}
			closer, isOpener := closeForOpen[c]
			switch {
			case isOpener:
				stack = append(stack, closer)
			default:
				mustBe := stack[len(stack)-1]
				if c != mustBe {
					isValid = false
				} else {
					stack = stack[:len(stack)-1]
				}
			}
		}
		if isValid {
			incompleteStacks = append(incompleteStacks, stack)
		}

		fmt.Printf("Line %s: %t\n", l, isValid)
	}

	sums := []int{}
	for _, stack := range incompleteStacks {
		subSum := 0
		for i := len(stack) - 1; i >= 0; i-- {
			pts := finalPoints[stack[i]]
			subSum = subSum*5 + pts
		}
		sums = append(sums, subSum)
		fmt.Printf("subsum: %d\n", subSum)
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

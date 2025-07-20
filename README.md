[![progress-banner](https://backend.codecrafters.io/progress/grep/0055a80a-5eb8-47ed-a4d5-415de8d1b85a)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is a starting point for Go solutions to the
["Build Your Own grep" Challenge](https://app.codecrafters.io/courses/grep/overview).

## INTRO

I am a student. A student at 26. I'm pushing 30 and I still have to worry about grades and homework. Jokes aside you do always have something to worry about. And I choose to worry about the fact that if I am stil a junior dev at 26. Maybe if I had not pursued a graduate education I would have been a mid level/senior dev by now.

So, what to do. How do I still get A's and continue/force myself to learn and level up as a software dev. Simple, throw money at it. I am a very frugal guy. I dont spend money unless I have to but I decided to throw my money on codecrafters. And it's working so far.

Quoting the legendary Primegan, doing big projects is the best way to learn. During my undergrad I did MERN and REST and easy things like these but now I have high standards for what I would like to do. I was working with a DB project for sometime but due to my phd commitments, excuses, and the complexity and having to navigate all that in 1 hr time I squeeze for myself. I made very slow progress. My professor is anything but an advocate for not getting stuck. I was stuck. This was another reason I decided to throw my money on codecrafter.

Now, onto the only thing that matters....

## WHAT DID I LEARN??

1. 1st time doing anything with GO. Beautiful langugage. Ken Thompson !!! Nuf said! . So, there is that. I am not unfamilar with modern powerful langauge like this because I've used Rust before and am in the Rust Cult :P. But yes beautiful language, easy to learn, succinct, modern. 

2. I approached this problem like the junior dev I am. The inital test cases were simple used loops and moved up the progression until I didn't.

3. Enter parsing. Recently i've not been able to find things not connected with compilers and or automata theory. My approach was like an DFA. Determintic approach. But parsing, and regular expression parsing with the intorduction of closure and kleen closure meant that a DFA was out. Enter NFA.

4. What's an NFA you ask? "NFA is magic", literlly something my professor who taught me Automata Theory said to help us better understand NFA.

5. Then i rewrote the entire thing. First things was expression parsing. 

6. PARSINGGG (PTSD from my compiler theory class where we had built an entire programming language using Java and Java Cup) is recursive magic. Non-Determinsitc problems like parsing are where recursions shine. 

    * I used something called a Recursive Descent Parsing(RDP). Basically, parse the regular expression based on prescedence. I had a though time figuring this out TBH. BUt I did. I am sad to say that under the timing constraint I used GPT's help with understand  RDP. I wrote the code myself (and i wrote it incorrect (knowingly) so that i knew why the algrithm was the way it was) but I needed help with the pesudo code. I guess this is how we read books now in the GPT era but me wanting to be a "purist" i dont enjoy it. But it's so efficient. 

    * With my RDP built i was in a much better place then when i was solving this whole problem in the previous way. Tree, graphs, DFS, man they just make everything so much earier. Oh and my RDP is a full binary tree (most of the time).

    * Next was to start building another parser. Using the RegExp built using RDP I needed to parse that tree with the checkString and see if there is indeed a grep match or not in the checkstring.

    ```
    checkString = "abcd"
    RegExpNode = Concat(Literal(b), Concat(Literal(c), Literal(d)))
    Res = checkParseTree(checkstring, regExpNode) //True since RegExp is a substring of checkString
    ```

    * Something I learned from spending too much time doing leetcode. If a probelm's solution does needs a lot of ifs then there exits a better general approach. But this is really just another way of saying "work hard so you can learn to work smart" - Primegan. This RDP apporach was so intuitive and the code was so clean when i started adding in parsing for concat, groups, anchors (^,$). 

    * Enter closures and kleen closures. I had ZERO idea that solving this problem requires such massive level of innovation and is such a huge engieering task. I mean look at [RE2](https://github.com/google/re2)

    * So, how you solve closures is important. I am doing backtracking with some herusitic based optimization cause that is all the time I have for this, initially. But RE2 uses NFA, and solves closure in linear time (O(checkString)). 
    "The more I learn, the more I realize how much I don't know." - Some genius scientist
    
    (learning things like these now a days is such a cheat code, before GPT i would find myself going down the rabbit hole for hours just looking at other people's solution and tyring to understand the gist of it, now GPT give me the summry and my precious time is saved and I can admire and get back to doing what I really need to do). 

    * I ran into a probelm pretty late
    "I am I saw I happy", check with " I saw happy"
    When i ran this i was getting true, even though it was false.
    I identified that my if and else conditions inside recursions were incorrect. It matched I and then saw and then happy and returned true but we need a contious match. The fix was very easy (backtracking with a heuristic). However, due to recursive nature the solution was not as easy as the many backtracking CP problems i've solved. So, I came up with counting the depth of the recursive calls (on my own :) ).Doing this i was able to apply the backtracking only at the start and not else where which was not needed as if anywhere during recursive parsing if i get false I need to go back all the way to the start and do and match the next char in checkString. 

## Analysing Complexity

1. Building the RegExp for the patter using RDP

Time Complexity: O(N), where N is the lenght of the pattern. This is because we go though each char in the pattern just once
Memory: O(N), Memory used depends on the nature of the tree created. Each chars are leafs to an operation with is the parent which is leaf to other operations
eg caa is concat(c,concat(a,a)), depth of tree is 3
or (c|a)*a is concat(repeat(alternate(c,a)), a), depth is 3
So, the memory used is based on the number of characters and not on the operations. Thus memory used is a binary tree with N nodes. O(N) maybe

2. Parsing RegExp with checkstring.

I dont know how, to anlyse the complexity for such a huge program tbh but ill try. I analyse the complexity of each sub cases. 

* Complexity for concat (Backtracking with one heuristic: skipping already checked strings)
Once the Abstract Syntax Tree(AST) is constructed with the pattern, the parsing is done in i think O(N) times, where N is the length of the checkString.
We go through each char in the string to see if there is a match with the current head left node. If match we recursive parse from next char position and the right node. If no match we continue to parse the next char with the first left node.

Eg: checkString: I am happy I saw
pattern: I saw
1. We get I match and then " " and then no match for "a" in "am". So we return back all the way skipp all checked chars " " and "a" and start from "m ...." and check each char to see if it matched with first leftChild that is I in that case. 
So, O(N)

* Complexity for literal, digit, alphaNum
O(N) or O(1) depends 
If we need to check literal match in concat we need contiguous match. So we cannot skip anything so, O(1). Else O(N) iterated over all till we find what we want. 
O(N) for checkString: "abcd1" pattern:\d

* Complexity for anchor start
O(1) when this is encountered the checkString should be in position 1

* Complexity for anchor end
O(1) when this is encountered the checkString should be in end position

* Complexity for repeat (Greedy Matching with Backtracking)
Best Case: O(n) eg checkString aaaaab, pattern a*b
Worst Case: O(2^n) eg checkString aaaab, pattern a*ab

* Complexity for alternate ( I dont know how to do this)
Normal case:
Eg: checkString abcd pattern ab(c|e)d check ab then check c or e then d
Complex Case:
eg: for pattern that alternates with repeats and other things i dont know prolly 2^n

* Complexity for Wildcard
O(1) just skip and match the rest
if .* then exponential again

## What can i do better?

1. Write better code, obv. The code i write today is horse shit tomorrow and that is by nature.

2. Write better error message.

3. Write better comments.
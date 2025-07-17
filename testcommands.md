### State # (ALternation)
echo -n "a cat" | ./your_program.sh -E "a (cat|dog)"
echo -n "a cow" | ./your_program.sh -E "a (cat|dog)"
echo -n "cat dog cow " | ./your_program.sh -E "(cat |dog |cow)+$" //true
echo -n "I see 1 cat, 2 dogs and 3 cows" | ./your_program.sh -E "^I see (\d (cat|dog|cow)s?(, | and )?)+$" //true

echo -n "1 cat, 2 dogs and 3 cows" | ./your_program.sh -E "(\d (cat|dog|cow)s?(, | and )?)+$" //true
echo -n "cat, dogs and cows" | ./your_program.sh -E "((cat|dog|cow)s?(, | and )?)+$" //true
echo -n "c, d and c" | ./your_program.sh -E "((c|d)(, | and )?)+$" //true
echo -n "c , c " | ./your_program.sh -E "(c |, )+$" //true


### Stage # (WildCard)
echo -n "cat" | ./your_program.sh -E "c.t" //true
echo -n "car" | ./your_program.sh -E "c.t" //false
echo -n "goøö0Ogol" | ./your_program.sh -E "g.+gol" //true
echo -n "gol" | ./your_program.sh -E "g.+gol" //false

### Stage # (Match zero or more times)
echo -n "dogs" | ./your_program.sh -E "dogs?" //true
echo -n "dog" | ./your_program.sh -E "dogs?" //true
echo -n "doga" | ./your_program.sh -E "dogs?a" //true
echo -n "dogssa" | ./your_program.sh -E "dogs?a" //true
echo -n "a" | ./your_program.sh -E "s?" //true // TODO: does not handle this

echo -n "cat" | ./your_program.sh -E "ca?t" //true
echo -n "act" | ./your_program.sh -E "ca?t" //true
echo -n "dog" | ./your_program.sh -E "ca?t" //false
echo -n "cag" | ./your_program.sh -E "ca?t" //false


### Stage #FZ7 (Match one or more times)
echo -n "cat" | ./your_program.sh -E "ca+t" //true
echo -n "caaats" | ./your_program.sh -E "ca+at" //true
echo -n "act" | ./your_program.sh -E "ca+t" //false
echo -n "ca" | ./your_program.sh -E "ca+t" //false

echo -n "acach" | ./your_program.sh -E "a(ca|ch)s+" //false
echo -n "acach" | ./your_program.sh -E "a(ca|ch)+" //true
echo -n "acack" | ./your_program.sh -E "(cj|ch)+" //false

### Stage #AO7 (End of string anchor)
echo -n "cat" | ./your_program.sh -E "cat$"
echo -n "cats" | ./your_program.sh -E "cat$"

### Stage #RR8 (Start of string anchor)
echo -n "log" | ./your_program.sh -E "^log"
echo -n "slog" | ./your_program.sh -E "^log"

### Stage #SH9 (Combining Character Classes)
echo -n "sally has 3 apples" | ./your_program.sh -E "\d apple" //true
echo -n "sally has 1 orange" | ./your_program.sh -E "\d apple" //false
echo -n "sally has 124 apples" | ./your_program.sh -E "\d\d\d apples" //true
echo -n "sally has 12 apples" | ./your_program.sh -E "\d\\d\\d apples" //false 
echo -n "sally has 3 dogs" | ./your_program.sh -E "\d \w\w\ws" //true
echo -n "sally has 4 dogs" | ./your_program.sh -E "\d \w\w\ws" //true
echo -n "sally has 1 dog" | ./your_program.sh -E "\d \w\w\ws" //false

### Stage #RK3 (Negative Character Groups)
echo -n "apple" | ./your_program.sh -E "[^xyz]" //true
echo -n "banana" | ./your_program.sh -E "[^anb]" //false
echo -n "orange" | ./your_program.sh -E "[^opq]" //true

### Stage #TL6 (Positive Character Groups)
echo -n "a" | ./your_program.sh -E "[abcd]"
echo -n "efgh" | ./your_program.sh -E "[abcd]"

### Stage #MR9 (Match alphanumeric characters)
echo -n "word" | ./your_program.sh -E "\w"
echo -n "$!?" | ./your_program.sh -E "\w"

### Stage #OQ2 (Match digits)
echo -n "123" | ./your_program.sh -E "\d"
echo -n "apple" | ./your_program.sh -E "\d"

### Stage #CQ2 (Match a literal character)
echo -n "dog" | ./your_program.sh -E "d"
echo -n "dog" | ./your_program.sh -E "f"

#!/bin/sh

#generate pdf files based on Polygon problem directory

source="$1" # Polygon package path
dest="$2" # problem path

workdir="$source/statement-sections/english"

cp scripts/gen_statement/gen_pdf.sh $workdir || exit -1
cp scripts/gen_statement/main.tex $workdir || exit -1
cp scripts/gen_statement/styles.sty $workdir || exit -1
cp scripts/gen_statement/logo.png $workdir || exit -1
cp scripts/gen_statement/replace_deprecated.sh $workdir || exit -1

cd $workdir

# ICPC type doesnt have scoring section, so this is psuedo-file
[ -f scoring.tex ] && printf "%smyheading{Chấm điểm}\n\n" "\\" | cat - scoring.tex 2>/dev/null > tmp && mv tmp scoring.tex
touch scoring.tex

[ -f output.tex ] && printf "%smyheading{Đầu ra}\n\n" "\\" | cat - output.tex 2>/dev/null > tmp && mv tmp output.tex
touch output.tex

[ -f input.tex ] && printf "%smyheading{Đầu vào}\n\n" "\\" | cat - input.tex 2>/dev/null > tmp && mv tmp input.tex
touch input.tex

[ -f examples.tex ] && printf "%smyheading{Ví dụ}\n\n" "\\" | cat - examples.tex 2>/dev/null > tmp && mv tmp examples.tex
touch examples.tex

[ -f notes.tex ] && printf "%smyheading{Chú thích}\n\n" "\\" | cat - notes.tex 2>/dev/null > tmp && mv tmp notes.tex
touch notes.tex

# [ -f interaction.tex ] && printf "%smyheading{Interaction}\n\n" "\\" | cat - interaction.tex 2>/dev/null > tmp && mv tmp interaction.tex
touch interaction.tex

./gen_pdf.sh || exit -1

cp ./statement.pdf "$dest" || exit -1

exit 0
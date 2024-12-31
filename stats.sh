#! /bin/bash

print_stats () {
	echo $1
	find . -type f -name "*.$2" -exec wc -l {} + | sort
	echo ""
}

print_stats "HTML" tmpl
print_stats "CSS" css
print_stats "JAVASCRIPT" js
print_stats "SQL" sql
print_stats "SHELL" sh
print_stats "GOLANG" go

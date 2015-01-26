go test -v
echo "finished testing ..............................................."
gocov test | gocov-html > coverage_report.html

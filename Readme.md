compare Site
===========================

# Description
This script downloads the web pages and calcuates the Jaccard Distance between the source code.
All references to the hostname are removed before doing so.


# Usage
Start the webserver with the following command.
`go run cmd/main.go`

Then browse to `http://localhost:8081`.
In the web interface enter the complete URI including the schema.
If everything works correct, a table with the results should show up.


# Links
* [Jaccard index](https://en.wikipedia.org/wiki/Jaccard_index)

greet:
  echo "Hello $1, my old friend"

deps: greet
  echo "Deps"
  
deps_error: greet error
  echo "Deps error"

error:
  exit 1

javascript:
  #!/usr/bin/env node
  console.log("hello")

python:
  #!/usr/bin/env python3
  print("hello")

perl:
  #!/usr/bin/env perl
  print "hello\n";

ruby:
  #!/usr/bin/env ruby
  puts "hello"
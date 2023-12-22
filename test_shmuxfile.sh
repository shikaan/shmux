greet:
  echo "Hello $1, my old friend"

deps: greet
  echo "Deps"
  
deps_error: greet error
  echo "Deps error"

error:
  exit 1
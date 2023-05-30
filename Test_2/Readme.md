# Word Service API

This is an HTTP service that allows you to store words and retrieve the most frequent word with a given prefix.

## Endpoints

### Store a Word

Stores a word in the service.

- URL: `/service/word`
- Method: `POST`
- Parameters:
  - `word` (string): The word to be stored.
- Example:
  ```bash
  curl -X POST -d "word=animal" http://localhost:8545/service/word 

 


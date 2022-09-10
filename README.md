
# Rubiks Api
This is a small Api written in go to serve random algorithms for cubing

# Query a specific algorithm
There are over 100 algorithms stored in the database and all can be queried by their name
Example: [Querying JPerm](https://rest-api-z7cayewqka-uc.a.run.app/v1/algorithm/JPerm)

# Query by category
There are 4 categories available OLL, PLL, F2l, 2LOOK
Example [Querying PLL(https://rest-api-z7cayewqka-uc.a.run.app/v1/algorithmCategory/PLL)

# Random OLL
The Api also provides one random OLL a day from any OLL in the database. 
This refreshes at 00:00 GMT
Example [random OLL](https://rest-api-z7cayewqka-uc.a.run.app/v1/randomOLL/)

# Random PLL
The Api also provides one random PLL a day from any PLL in the database. 
This refreshes at 00:00 GMT
Example [random PLL](https://rest-api-z7cayewqka-uc.a.run.app/v1/randomOLL/)

## Features

- **Current Block Information**: Fetch and display the current block number of the Ethereum blockchain, allowing users to stay up-to-date with the latest block.

- **Address Subscription**: Users can subscribe to Ethereum addresses of interest. This feature enables the monitoring of transactions related to these addresses, enhancing the user's ability to track and manage blockchain activities.

- **Local Storage of Transactions**: Transactions related to subscribed addresses are fetched and stored locally. This design choice ensures quick access to transaction data and reduces reliance on external calls to the blockchain for historical data.

- **Retrieval of Transactions**: Allows users to retrieve a list of transactions associated with a subscribed address. It's important to note that this feature retrieves transaction data from local storage, rather than fetching it in real-time from the blockchain. This approach significantly improves response times and efficiency when accessing historical transaction data.

## Components

The project is structured into several key components, each responsible for a distinct aspect of the application:

- **Main Server**: The entry point of the application, setting up HTTP endpoints and initializing the blockchain parser and storage.

- **Blockchain Parser**: Interfaces with the Ethereum blockchain, fetching current blocks, and processing transactions for subscribed addresses.

- **Storage**: Manages the in-memory storage of transactions and subscribed addresses, ensuring quick data retrieval and efficient storage management.

- **Client**: Facilitates JSON-RPC communication with Ethereum nodes, allowing the application to interact with the blockchain for fetching blocks and transactions.

## Usage

1. **Starting the Server**: Run the main server to start listening for incoming requests on port 8080. This sets up endpoints for fetching the current block, subscribing to addresses, and retrieving transactions.

2. **Fetching the Current Block**: Access the `/currentBlock` endpoint to get the latest block number on the Ethereum blockchain.

3. **Subscribing to an Address**: Use the `/subscribe` endpoint with an address query parameter to subscribe to an Ethereum address. The application will then monitor transactions related to this address.

4. **Retrieving Transactions**: Access the `/transactions` endpoint with an address query parameter to retrieve a list of transactions associated with the subscribed address from local storage.

## Emphasis on Local Storage

The `GetTransactions` functionality emphasizes retrieving transaction data from local storage rather than making real-time calls to the Ethereum blockchain. This design decision ensures that users can quickly access historical transaction data without the latency associated with blockchain queries. It's an efficient solution for monitoring and analyzing blockchain activities over time.

## TO DO

- **Implementing Transaction Saving**: A critical feature yet to be implemented is the automatic fetching and saving of transactions from the Ethereum blockchain to the local storage. This involves monitoring the blockchain for new transactions related to subscribed addresses and storing these transactions locally for quick access. The implementation should handle:

  - **Fetching New Transactions**: Regularly check the Ethereum blockchain for new transactions associated with subscribed addresses.

  - **Processing Transactions**: Once new transactions are detected, process and format these transactions to match the local storage schema.

  - **Saving to Local Storage**: Efficiently save the processed transactions to local storage, ensuring they are readily accessible for retrieval through the `GetTransactions` endpoint.

This functionality is essential for keeping the local transaction database up-to-date and ensures that users have access to the latest transaction information without having to query the blockchain directly.

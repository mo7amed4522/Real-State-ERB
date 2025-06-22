# Chat API Example Queries

## Create Room
```graphql
mutation {
  createRoom(input: {
    name: "General",
    description: "General chat room",
    type: "group",
    participantIds: [2, 3]
  }) {
    id
    name
    participants { userID }
  }
}
```

## Send Message
```graphql
mutation {
  sendMessage(input: {
    roomId: "1",
    content: "Hello world!",
    messageType: "text"
  }) {
    id
    content
    createdAt
  }
}
```

## Reply to Message
```graphql
mutation {
  sendMessage(input: {
    roomId: "1",
    content: "Replying to this",
    messageType: "text",
    replyToId: "5"
  }) {
    id
    content
    replyTo { id content }
  }
}
```

## Add Emoji Reaction
```graphql
mutation {
  addReaction(input: {
    messageId: "1",
    emoji: "👍"
  }) {
    id
    emoji
  }
}
```

## Upload File
```graphql
mutation($file: Upload!) {
  uploadFile(input: {
    messageId: "1",
    file: $file
  }) {
    id
    fileName
    encryptedPath
  }
}
```

## Subscribe to Messages
```graphql
subscription {
  messageAdded(roomID: "1") {
    id
    content
    senderID
    createdAt
  }
}
``` 
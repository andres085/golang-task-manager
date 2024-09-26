package main

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const userIDContextKey = contextKey("userID")

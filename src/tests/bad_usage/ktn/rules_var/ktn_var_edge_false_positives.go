package rules_var

import "os"

// maxRetriesLiteral est un littéral jamais modifié (DEVRAIT être const VAR-005).
var maxRetriesLiteral int = 3

// validHTTPCode utilise l'initialisme HTTP (devrait être const VAR-005).
var validHTTPCode int = 200

// maxHTTPRetries utilise HTTP en majuscules (devrait être const VAR-005).
var maxHTTPRetries int = 5

// max_HTTP_Retries utilise underscore (ERREUR VAR-008).
var max_HTTP_Retries int = 5

// MAX_TIMEOUT est en ALL_CAPS (ERREUR VAR-009).
var MAX_TIMEOUT int = 30

// urlFromEnv lit une URL depuis l'environnement (OK car appel fonction).
var urlFromEnv string = os.Getenv("SERVICE_URL")

// computedValue est calculé via une fonction (OK car appel fonction).
var computedValue int = len("hello")

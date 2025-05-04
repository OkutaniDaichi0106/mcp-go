# Prompt for Generating Code Comments

You are writing a comment for a Go function or block of code. The comment must clearly describe:
- The purpose of the function or code block
- The expected inputs and their types
- The output(s) and their types
- Any side effects or important behavior
- Any error conditions or edge cases handled
- (If applicable) why a specific approach or algorithm is used

Write the comment in natural, professional English as if explaining to a peer Go developer. Use full sentences and proper punctuation.

Example Format:
~~~go
// FetchUserProfile queries the user profile by the given user ID.
// It returns a User object if found, or an error if the user is not present
// or the database query fails.
~~~

Now, generate a similar comment for the following function:
<INSERT GO CODE SNIPPET HERE>

// create a logged in home page view
// at the very top of the page, provide the user the ability to click to go to their profile
// at the top, provide a small amount of direction for users
// show large icons for each steaming service
// that when clicked take you to the FWIP start menu

// create a logged in FWIP start menu
// at the very top of the page, provide the user the ability to click to go to their profile
// at the top of the page, create a filter to allow users to toggle what results are delivered to them (tv, movies, or all content)
// Provide a menu with the following options:
// include a note that changes submitted via the menu may take some time to appear live on the site
// 1- Allow the user to update an existing title in the streaming service
// by filling out a form in a <dialog> window
// user must include the title they're updating
// a "steaming service" dropdown should be pre-populated with the streaming service they've selected
// 2- Allow the user to move or delete a title from the streaming service
// by filling out a form in a <dialog> window
// user must indicate if they are moving or deleting a title
// if deleting a title, user must include a reason why they're removing the title
// user must include the title they're updating or removing
// a "steaming service" dropdown should be pre-populated with the streaming service they've selected
// 3 - Add a new title
// by filling out a form in a <dialog> window
// // user must include the title they're adding
// a "steaming service" dropdown should be pre-populated with the streaming service they've selected
// all changes will be pushed to a "moderation queue" API that will then be delivered to a moderation tool
// Under these smaller buttons should be a larger button that says, “Fwip!”
// when clicked this button will take you to a random title card that meets the criteria selected

// the moderation tool should
// show the admin any changes requested by a user, noting the source
// 1- from the "edit dialog" source
// 2 - from the "add dialog" source
// 3 - from the "move dialog" source
// 4 - from the "delete dialog" source
// the ability to modify any user input in another dialog box
// a button to "Approve Changes"
// for the edit, add, and move source, this will change the API in accordance with the submission
// for the "delete" source this approval will DELETE the content from the API
// The "Approve Changes" button will push the changes into the main API

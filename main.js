// state

const BASE_API = `https://kayabuki.neon-sole.ts.net/titles`;
// this API should contain the basic API URL structure
// extension API for service selection will include ?service="serviceId"
// netflix, disney, amazon, max, paramount, etc
// create API for IMDB ID of the title
// moderation API

// title (which has an ID)
// -> which is either a movie or a tv show
// -> contains episodes (if movie, episode count is 0)

// home API pull and render

const getTitles = async () => {
  let fetchAPI = await fetch(BASE_API);
  let transformerJson = await fetchAPI.json();
  return transformerJson;
};

// create a logged in home page view
// at the very top of the page, provide the user the ability to click to go to their profile
// at the top, provide a small amount of direction for users
// show large icons for each steaming service
// that when clicked take you to the FWIP start menu

// start menu pull and render

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

// FWIP pull and render

// create the FWIP page
// see a random title ID from the selected streaming service, with the filter applied if you chose
// the title "card" will contain the following information:
//  the image for each title
//  the title of each title
//  The year the title was released
//  The runtime of the content in minutes
//  If a TV show, average episode length
//  The Rotten Tomatoes rating, if any
//  The director
//  Up to 10 of the top billed cast
//  A brief description of the film (no more than 500 characters)
// Under the title card will be three buttons:
//  Seen It
//      Clicking this button will change the watched value to true
//      The information will need to be pushed to the API for use in the profile page
//      In the profile, use this info to generate a list of titles the user has seen under each streaming service
//  Star Icon
//      Clicking this button will change the importance value to true
//      Having an importance value should eventually weight answers so when we ask the app to generate a watch list in the future, it will have a value for the important titles higher than titles that haven’t been seen.
//  Haven’t seen It
//      Clicking this button will change the watched value to false.
//      False titles will not appear on the profile
// Immediately under the card, place an “Edit Title” button to generate a <dialog> to allow the user to edit the title ID.

// Fwip! Action Requirements
// Indicate that they have seen the title
//      Place a true value on the title ID, gained from the hash of the page
//      Locate the hash identifier from the URL and have the user API store the positive value
// Indicate that they have not seen the title
//      Place a false value on the title ID, gained from the hash of the page
//      Locate the hash identifier from the URL and have the API store the false value
// Indicate that they would like to watch the title
//      Set the importance of this value to 1 and anything without importance to 0
//      Store this value in the user API
//      This value will be used to generate a weight-calculated watch list

// the moderation tool should
// show the admin any changes requested by a user, noting the source
// 1- from the "edit dialog" source
// 2 - from the "add dialog" source
// 3 - from the "move dialog" source
// 4 - from the "delete dialog" source
// the ability to modify any user input in another dialog box
// a button to "Approve Changes"
// for the edit, add, and move source, this will change the main API in accordance with the submission
// for the "delete" source this approval will DELETE the content from the API


function check(formdata) {
  console.log(formdata.fname.value)
  console.log(formdata)
  alert(formdata)
}


async function getqry () { // stand get using qeury string

  // post request using options object
  console.time("get response took");
  let response = await fetch('http://localhost:4001/crudops?Id=612ba818a6f253918807f123')
  let convert = await response.json ()
  console.timeEnd("get response took");
  console.log("get request completed:")
  console.log(convert)

}


async function getbody (selectedartists) { // same as get but using req body 

  //let data2snd ={id:['612ba818a6f253918807f123']}
  
  let data2snd ={id:'612ba818a6f253918807f123', selected: selectedartists
  }

  let options = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(data2snd),
  }

  // post request using options object
  console.time("getbody response took");
  let response = await fetch('http://localhost:4001/getbody', options)
  let convert = await response.json ()
  console.timeEnd("getbody response took");
  console.log("get request using body completed:")
  console.log(convert)

}


async function chkclashes (selectedartists) {

  let data2snd ={id:'612ba818a6f253918807f123', selected: selectedartists
  }

  let options = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(data2snd),
  }

  // post request using options object
  console.time("chkclashes response took");
  let response = await fetch('http://localhost:4001/clashes', options)
  let convert = await response.json ()
  console.timeEnd("chkclashes response took");
  console.log("Clash check completed, clashing artists are:")
  console.log(convert)// js api response
  console.log(convert.Clashes)// add . Clashes for go api json format

}


// put request to update doc using id
async function put () {

  // data to be sent in request to api includinb id of mongodb document
  let data = {Id:'6155d68e7f8b8cb0900f3f1f', Updatetype: '$set',
              Key2updt: 'Data1', Value2updt: 'Fuk'
  }

  // object containing fetch options for the put request
  let options = {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(data),
  }

  // post request using options object
  console.time("put response took");
  let response = await fetch('http://localhost:4001/crudops', options)
  let convert = await response.json ()
  console.timeEnd("put response took");
  console.log("Put completed updated document is:")
  console.log(convert)

}

//put()

// post request

async function post (data2pst, endpoint) { // takes data to put in body and endpoint to post to

  let data2snd = {
    Data1: data2pst,
  }

  let options = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(data2snd),
  }

  // post request using options object
  console.time("post response took");
  let response = await fetch('http://localhost:4001/'+endpoint, options)
  let convert = await response.json ()
  console.timeEnd("post response took");
  console.log("Post completed created document is:")
  console.log(convert)
  return convert

}

//post()

//delete request

async function deleted () {

  let data2del = {
    1: 'hiya',
    2: 'cunt'
  }

  let options = {
  method: 'DELETE',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(data2del),
  }


  // post request using options object
  console.time("deleted response took");
  let response = await fetch('http://localhost:4001/crudops', options)
  let convert = await response.json ()
  console.timeEnd("deleted response took");
  console.log(convert)

}

//deleted()





async function signin(){
  console.log("called func")
  const username = document.getElementById('username').value;// get username and password
  const text = document.getElementById('pass').value;
  //console.log(username, password)
  let ws = new WebSocket("ws://86.148.246.175:4444/websocket")

  ws.addEventListener("message", function (e) {
    let data = JSON.parse(e.data);
    let chatContent = `<p><strong>${data.username}</strong>: ${data.text}</p>`;
    console.log(data)
    //room.append(chatContent);
    //room.scrollTop = room.scrollHeight; // Auto scroll to the bottom
  });
  

  ws.onopen = () => ws.send(
    JSON.stringify({
      username: username.value,
      text: text.value,
    })
  );
  
  

/*   // post request to check sign in details
  console.time("signin response took");
  let response = await fetch(signinurl, options)// post signin details with id in url
  let convert = await response.json ()
  console.timeEnd("signin response took");
  console.log("User password check complete:")
  console.log(response.status)
  if ( response.status === 200){ // if password is a match redirect to profile page
      console.log("pword correct cookie is", "authjwt="+convert.token)
      document.cookie = "authjwt="+convert.token; // store data in cookie needs to be http only do later!
      //window.location.href = "http://www.w3schools.com";
  } else if ( response.status === 401){// handle not a match
      alert("userame or password was incorrect please try again")
  } */

}




async function uploadpic (picnum, imgreplaceid) {
  var input = document.getElementById(picnum)// get file thats been given
  var data = new FormData() // create empty form data cos api reads form for the file..coulda done this.form but gave errors
  data.append('file', input.files[0]) // add the file to the form
  data.append('user', 'hubot')
 
  let options = {
  method: 'POST',
  body: data, // send the form with the file and user data
  }

  username = "betterbepub"
  signinurl = 'http://localhost:4001/postfile/' + username + "/" + picnum// get username from cookie?

  // post the json
 // console.time("deleted response took");
  let response = await fetch(signinurl, options)
  let convert = await response.json ()
  //console.timeEnd("deleted response took");
  console.log(convert.Imgaddy)
  newimageurl = convert.Imgaddy // get posted img address and change profile picture
   pic2replace = picnum+"1"
  document.getElementById(pic2replace).src = newimageurl;

}

async function signup () { // sends userbame password from input(need!) server creates bucket in username and stores details on mongo
  
  const username = document.getElementById('signupusername').value;// get username password and email
  const password = document.getElementById('signuppass').value;
  const email = document.getElementById('signupemail').value;
    
  let signupdata = {
    Username: username,
    Password: password,
    Email: email
  }

  let options = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(signupdata),
  }
  
  signinurl = 'http://localhost:4001/signup/' + username

  // post request using options object
  console.time("signup response took");
  let response = await fetch(signinurl, options)// need to make string with id at end from input
  let convert = await response.json ()
  console.timeEnd("signup response took");
  console.log("User signup completed created document is:")
  console.log(convert)

}

//signup()

function getcookie (name) { // trims document.cookie string returns value
   //document.cookie = "username=John Doe";
   //let x = document.cookie;
  console.log("cookiez is", x) 
  const cookieValue = document.cookie
  .split('; ')
  .find(row => row.startsWith(name))
  .split('=')[1];
  return cookieValue
}

function chkauth (){//chk if jwt is valied for page used by each html on every load
    jwt = getcookie("authjwt")
    response = post(jwt, "chkauth")
    if ( response.status === 401){// handle not a match
      alert("you not authorized to view this page please sign in")
      window.location.href = "signin.html";     
     } else if ( response.status === 200){ // if password is a match redirect to profile page
      console.log("authorized by api")
     }
}

function triggerclik(name){

// clicking image triggers upload button click
var myButton = document.getElementById(name);
    myButton.click()
}

//document.getElementById("inputpic1").onchange = function() {
//    sayhi()
//    button2 = document.getElementById("uploader");
//    button2.click()
//};

function sayhi () {
    //var input = document.querySelector('input[type="file"]')// get tile thats been given
  //consle.log("file input format is")
 // consle.log(input.files[0])
    alert("hi")
}

//sayhi()



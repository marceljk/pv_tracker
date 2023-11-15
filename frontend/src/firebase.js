import { initializeApp } from 'firebase/app'
import { getDatabase, ref as dbRef } from 'firebase/database'
// ... other firebase imports

export const firebaseApp = initializeApp({
    apiKey: "AIzaSyDyCp-VqrhFqXbgjgL6_-SM86ljEwH0Wsk",
    authDomain: "poc-database-da139.firebaseapp.com",
    databaseURL: "https://poc-database-da139-default-rtdb.europe-west1.firebasedatabase.app",
    projectId: "poc-database-da139",
    storageBucket: "poc-database-da139.appspot.com",
    messagingSenderId: "326067718177",
    appId: "1:326067718177:web:5a8f0c2cd7682a785d3b60",
    measurementId: "G-5BCYEL3RTK"
})

// used for the databas refs
const db = getDatabase(firebaseApp)

// here we can export reusable database references
export const history = dbRef(db, 'history')
export const live = dbRef(db, 'live')

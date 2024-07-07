/*
  The data-refresher script is use to trigger refreshes to mempool and 
  bundle data.
*/
let intervalId = undefined

const startInterval = () => {
  intervalId = setInterval(() => {
    console.log("This code runs every 10 seconds");
  }, 10000);
};

const main = () => {
  try {
    if (intervalId === undefined) {
      startInterval();
      console.log("Interval started");
    } else {
      console.log("Interval is already running");
    }
  } catch (err) {
    console.error(err);
  }
};

// Automatically start the interval on page load
main();

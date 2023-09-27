const express = require('express');
const request = require('request');
const FormData = require('form-data');
const fs = require('fs');
const _ = require("lodash");

const { initializeApp, cert } = require('firebase-admin/app');
const { getDatabase } = require('firebase-admin/database');

const serviceAccount = require('./serviceAccountKey.json');

initializeApp({
  credential: cert(serviceAccount),
  databaseURL: "https://poc-database-da139-default-rtdb.europe-west1.firebasedatabase.app",
});

const db = getDatabase();
const dbHistoryRef = db.ref("history");
const dbLiveRef = db.ref("live");
const dbDailyForecastRef = db.ref("dailyForecast");
const dbWeeklyForecastRef = db.ref("weeklyForecast");
const dbDailySumRef = db.ref("dailySum");

const router = express.Router();

let cookie = "";

const login = () => {
  const form = new FormData();
  form.append("username", "user1");
  form.append("password", "ASrIJY");
  form.submit('http://varta130104162/cgi/login', (err, res) => {
    if (err) setTimeout(login, 1000 * 30);
    cookie = res.headers['set-cookie'];
  })
}

login();

router.get('/', function (req, res) {
  res.send("API !");
});

const getPVData = async () => {
  return await new Promise((resolve, rejects) => {
    request(
      {
        url: 'http://varta130104162/cgi/data',
        headers: {
          cookie: cookie,
        }
      },
      (error, response, body) => {
        if (error || response.statusCode !== 200 || body == undefined) {
          rejects(response);
        }

        try {
          const data = JSON.parse(body)
          resolve({
            batteryPercent: data.pulse.procImg.soc_pct,
            gridPower: data.pulse.procImg.gridPower_W,
            pvPower: data.common.power_W,
            powerConsumption: Math.round(data.common.power_W - data.pulse.procImg.gridPower_W - data.pulse.procImg.activePowerAc_W),
            batteryPower: Math.round(data.pulse.procImg.activePowerAc_W)
          });
        } catch (err) {
          console.log("getPVData() err:", err);
          console.log("getPVData() body:", body);
          rejects(err);
        }
      }
    )
  });
}

router.get('/data', async (req, res) => {
  const result = await getPVData();
  if (result.statusCode !== undefined) res.statusCode(result.statusCode).json(response);
  res.json(result);
});

router.get('/forecast', (req, res) => {
  const rawData = fs.readFileSync('forecast.json');
  const forecast = JSON.parse(rawData);

  let arr = [...Array(18).keys()].map(x => x * 1);

  const startEntry = arr[0];
  //Shift the data until the first record is not 0 (max 16 hours)
  for (let i = startEntry; i < (startEntry + 32); i++) {
    if (forecast.pv_estimate[i] == 0) {
      arr = arr.map(x => x + 1);
    } else {
      break;
    }
  }

  if (Date.parse(forecast.period_end[1]) < Date.now()) {
    arr = arr.map(x => x + 2);
  } else if (Date.parse(forecast.period_end[0]) < Date.now()) {
    arr = arr.map(x => x + 1);
  }

  let response = {
    pv_estimate: [],
    period_end: [],
    period: "",
  };
  arr.forEach((index) => {
    if (Math.ceil(forecast.pv_estimate[index]) > 0) {
      response.pv_estimate.push(forecast.pv_estimate[index]);
      response.period_end.push(forecast.period_end[index]);
    }
  });
  response.period = forecast.period;

  res.send(response);
});

router.get('/forecast/daily', (req, res) => {
  const rawData = fs.readFileSync('dailyForecast.json');
  const dailyForecast = JSON.parse(rawData);
  res.send(dailyForecast);
});

const getForecast = () => {
  const options = {
    method: 'GET',
    url: 'https://api.solcast.com.au/world_pv_power/forecasts?latitude=49.241698&longitude=9.058941&capacity=5&tilt=50&azimuth=180&install_date=2021-09-01&hours=168&api_key=EOSLqAy1Qml0rRTdQvu4_9E10YZNPgBo',
    //url: 'http://localhost:3030/forecast',
    headers: {
      accept: "application/json"
    },
  }
  request(options, async (err, res) => {
    if (!err && res.statusCode == 200) {
      const data = await JSON.parse(res.body);

      let temp = {};
      data.forecasts.forEach((x) => {
        let period_date = new Date(x.period_end).setHours(0, 0);
        if (temp[period_date] == undefined) {
          temp[period_date] = x.pv_estimate;
        } else {
          temp[period_date] += x.pv_estimate;
        }
      });

      Object.keys(temp).forEach(x => {
        temp[x] = (Math.round((temp[x] / 2) * 1) / 1);
      });

      let dailyForecast = {
        dailyForecast: {}
      };
      Object.keys(temp).forEach(x => {
        const day = new Date(Number(x)).getDay();
        dailyForecast.dailyForecast[x] = {
          estimate: temp[x],
          day: getDay[day]
        };
      }
      );

      fs.writeFile("dailyForecast.json", JSON.stringify(dailyForecast), function (err) {
        if (err) {
          console.log("err:", err);
        }
      });

      const forecast = {
        pv_estimate: data.forecasts.map(x => x.pv_estimate),
        period_end: data.forecasts.map(x => x.period_end),
        period: data.forecasts[0].period,
      };

      try {
        dbDailyForecastRef.set(dailyForecast.dailyForecast);
        dbWeeklyForecastRef.set(forecast);
      } catch { }

      fs.writeFile("forecast.json", JSON.stringify(forecast), function (err) {
        if (err) {
          console.log(err);
        }
      });
      console.log("updated forecast", new Date());
    } else {
      console.log(err);
    }

  });
};

const storeDataInFirebase = async () => {
  try {
    const res = await getPVData();
    counter++;
    if (counter > 9) {
      const now = new Date().toISOString().substring(0, 19);
      dbHistoryRef.child(now).set(res);
      counter = 0;
    }
    dbLiveRef.set(res);
  } catch (err) {
    console.log(err);
  }
}

const sortData = {};
const addDailySum = () => {
  dbHistoryRef.orderByKey().on('child_added', (data) => {
    const key = data.key.substring(0, 11);
    const value = data.val();
    let gridPowerOut = 0;
    let gridPowerIn = 0;
    const gridPower = value.gridPower;
    if (value.gridPower > 0) {
      gridPowerIn = gridPower;
    } else if (value.gridPower < 0) {
      gridPowerOut = gridPower;
    }

    if (sortData[key]) {
      sortData[key] = {
        batteryPercent: sortData[key].batteryPercent + value.batteryPercent,
        batteryPower: sortData[key].batteryPower + value.batteryPower,
        gridPowerIn: sortData[key].gridPowerIn + gridPowerIn,
        gridPowerOut: sortData[key].gridPowerOut + gridPowerOut,
        powerConsumption: sortData[key].powerConsumption + value.powerConsumption,
        pvPower: sortData[key].pvPower + value.pvPower,
        count: sortData[key].count + 1,
      }
    } else {
      sortData[key] = {
        batteryPercent: value.batteryPercent,
        batteryPower: value.batteryPower,
        gridPowerIn: gridPowerIn,
        gridPowerOut: gridPowerOut,
        powerConsumption: value.powerConsumption,
        pvPower: value.pvPower,
        count: 1,
      }
    }
    saveDailySum();
  });
};

const saveDailySum = _.debounce(() => {
  Object.entries(sortData).forEach((obj) => {
    const key = obj[0];
    const value = obj[1];
    dbDailySumRef.child(key).update(value);
    const today = new Date().toISOString().substring(0, 11);
    if (key < today) {
      delete sortData[key];
    }
}

const getDay = ['Sonntag', 'Montag', 'Dienstag', 'Mittwoch', 'Donnerstag', 'Freitag', 'Samstag'];
getForecast();
setInterval(getForecast, 1000 * 60 * 30);
setInterval(login, 1000 * 60 * 60 * 24);
setInterval(storeDataInFirebase, 1000 * 30);
const PORT = process.env.PORT || 3030;

const app = express();
app.use(router);

app.listen(PORT, () => console.log(`listening on ${PORT}`));

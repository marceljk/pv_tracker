import express from 'express';
import request from 'request';
import FormData from 'form-data';
import fs from 'fs';

const router = express.Router();

let cookie = "";

const login = () => {
    const form = new FormData();
    form.append("username", "user1");
    form.append("password", "ASrIJY");
    form.submit('http://varta130104162/cgi/login', (err, res) => {
        cookie = res.headers['set-cookie'];
    })
}

login();

router.get('/', function (req, res) {
    res.send("API !");
});

router.get('/data', (req, res) => {
    request(
        {
            url: 'http://varta130104162/cgi/data',
            headers: {
                cookie: cookie,
            }
        },
        (error, response, body) => {
            if (error || response.statusCode !== 200) {
                return res.status(response.statusCode).json(response);
            }

            const data = JSON.parse(body)
            res.json({
                batteryPercent: data.pulse.procImg.soc_pct,
                gridPower: data.pulse.procImg.gridPower_W,
                pvPower: data.common.power_W,
                powerConsumption: Math.round(data.common.power_W - data.pulse.procImg.gridPower_W - data.pulse.procImg.activePowerAc_W),
                batteryPower: Math.round(data.pulse.procImg.activePowerAc_W)
            });
        }
    )
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
        if (err == null) {
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
                    console.log(err);
                }
            });

            const forecast = {
                pv_estimate: data.forecasts.map(x => x.pv_estimate),
                period_end: data.forecasts.map(x => x.period_end),
                period: data.forecasts[0].period,
            };
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

const getDay = ['Sonntag', 'Montag', 'Dienstag', 'Mittwoch', 'Donnerstag', 'Freitag', 'Samstag'];
getForecast();
setInterval(getForecast, 1000 * 60 * 30);
setInterval(login, 1000 * 60 * 60 * 24);
const PORT = process.env.PORT || 3030;

const app = express();
app.use(router);

app.listen(PORT, () => console.log(`listening on ${PORT}`));

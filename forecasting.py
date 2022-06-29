import sqlite3
import requests
# import pandas as pd
from datetime import date, datetime

def get_json_from_AP():
  apName = "apa05-0mg"
  eduroam     = "alias(ap.{ap}.ssid.eduroam,\"eduroam\")".format(ap = apName)
  lrz         = "alias(ap.{ap}.ssid.lrz,\"lrz\")".format(ap = apName)
  mwnEvents   = "alias(ap.{ap}.ssid.mwn-events,\"mwn-events\")".format(ap = apName)
  BayernWLAN  = "alias(ap.{ap}.ssid.@BayernWLAN,\"@BayernWLAN\")".format(ap = apName)
  other       = "alias(ap.{ap}.ssid.other,\"other\")".format(ap = apName)
  target      = "cactiStyle(group({eduroam},{lrz},{mwn_events},{bayern},{other}))".format(eduroam = eduroam, lrz = lrz, mwn_events = mwnEvents, bayern = BayernWLAN, other = other)
  url = "http://graphite-kom.srv.lrz.de/render?target={target}&format=json&from=-30days".format(target = target)
  jsonResp = requests.get(url).json()
  return jsonResp

def process_json(jsonResp):
  total = []
  for datapoint in jsonResp[0]['datapoints']:
    connDev = datapoint[0]
    ts = datapoint[1]
    tm = datetime.utcfromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S')
    if connDev is None:
      total.append([ts, 0])
    else:
      total.append([ts, connDev])

  for jsonEntry in jsonResp[1:]:
    datapoints = jsonEntry['datapoints']
    for idx, dpArr in enumerate(datapoints):
      ts = dpArr[1]
      tm = datetime.utcfromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S')
      connDevices = dpArr[0]
      currDev = total[idx][1]
      if connDevices is not None:
        total[idx] = [tm, currDev + connDevices]
      else:
        total[idx] = [tm, currDev]
  return total

json = get_json_from_AP()
processed = process_json(json)
for time, val in processed:
  print(time, val)
# print(processed)

def forecast():
  df1 = pd.DataFrame(total)
  df1.columns = ['ds', 'y']
  print(df1)

  m = Prophet()
  df1['floor'] = 0.0
  m.fit(df1)
  future = m.make_future_dataframe(periods=15)
  future['floor'] = 0.0
  future.tail()

  forecast = m.predict(future)
  forecast[['ds', 'yhat', 'yhat_lower', 'yhat_upper']].tail()
  fig1 = m.plot(forecast)
  fig2 = m.plot_components(forecast)

def process_forecasted_data(forecastedData):
  col_list = ["ds", "trend"]
  df = pd.read_csv(forecastedData, usecols=col_list)

  hourlyAvgs = []
  prevHour = -1
  trendTot = 0
  cnt = 0
  for _, row in df.iterrows():
    dateObj = pd.to_datetime(row['ds'])
    hour = dateObj.hour
    day = dateObj.day
    trend = row['trend']
    if prevHour < 0:
      prevHour = hour
      trendTot = trend
      cnt = 1
      continue
    if hour != prevHour:
      # calculate trend avg. for prevHour
      hourlyAvgs.append((day, prevHour, trendTot / cnt))
      prevHour = hour
      trendTot = trend
      cnt = 1
    else:
      trendTot += trend
      cnt += 1
    minutes = dateObj.minute
    if hour == 23 and minutes == 45:
      hourlyAvgs.append((day, prevHour, trendTot / cnt))
      prevHour = -1

  return hourlyAvgs

def write_to_DB(forecastData, name):
  con = sqlite3.connect('./data/sqlite/heatmap.db')
  cur = con.cursor()
  for day, hr, avg in forecastData:
    column = 'T' + str(hr)
    stmt = 'UPDATE future SET ' + column + ' = ? WHERE Day = ? AND AP_Name = ?'
    cur.execute(stmt, (avg, day, name))
  con.commit()

fdata = [(0, 1, 44), (0, 2, 88)]
# write_to_DB(fdata, "apa08-1w4")

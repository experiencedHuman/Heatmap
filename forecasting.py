import pandas as pd
from prophet import Prophet
import requests
from datetime import datetime

apName = "apa05-0mg"
eduroam     = "alias(ap.{ap}.ssid.eduroam,\"eduroam\")".format(ap = apName)
lrz         = "alias(ap.{ap}.ssid.lrz,\"lrz\")".format(ap = apName)
mwnEvents   = "alias(ap.{ap}.ssid.mwn-events,\"mwn-events\")".format(ap = apName)
BayernWLAN  = "alias(ap.{ap}.ssid.@BayernWLAN,\"@BayernWLAN\")".format(ap = apName)
other       = "alias(ap.{ap}.ssid.other,\"other\")".format(ap = apName)
target      = "cactiStyle(group({eduroam},{lrz},{mwn_events},{bayern},{other}))".format(eduroam = eduroam, lrz = lrz, mwn_events = mwnEvents, bayern = BayernWLAN, other = other)
url = "http://graphite-kom.srv.lrz.de/render?target={target}&format=json&from=-30days".format(target = target)
jsonResp = requests.get(url).json()


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

#print(total)
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
print("Plotting finished")
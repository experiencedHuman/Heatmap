import sqlite3

def get_ap_names():
  conn = sqlite3.connect("./data/sqlite/heatmap.db")
  cur = conn.cursor()
  cur.execute("SELECT DISTINCT Name FROM apstat")
  # rows = cur.fetchall()

  for row in cur:
    print(row[0])

get_ap_names()
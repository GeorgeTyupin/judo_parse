import pandas as pd
import json


# Путь к файлу
file_path = 'Соревнования.xlsx'

sheet_names = pd.ExcelFile(file_path).sheet_names

to_json = {}

def find_len_tables(df):
    len_tables = []

    arr = []
    for i in df.iloc[0]:
        if i == "|": 
            arr.append(1)
        elif i == "end":
            arr.append(2)
        else:
            arr.append(0)

    i = 0
    n = len(arr)
    while i < n:
        if arr[i] == 0:
            start = i
            while i < n and arr[i] == 0:
                i += 1
            len_tables.append(i - start)
            if i < n and arr[i] == 2:
                break
        else:
            i += 1
    return len_tables

for name in sheet_names:
    # Загрузка таблицы, пропуская первые 2 строки
    df = pd.read_excel(file_path, skiprows=2, sheet_name=name)
    tournaments = []
    left = 1
    for i, len_current_df in enumerate(find_len_tables(df)):
        right = left + len_current_df
        current_df = df.iloc[:, left:right]
        headlines =current_df.iloc[5]
        left = right + 1
        tournaments.append({
            "name" : current_df.iloc[0, 0],
            "description" : current_df.iloc[1, 0],
            "date" : current_df.iloc[2, 0],
            "gender" : current_df.iloc[3, 0], 
            "weight_categories" : {}
        })

        current_weight_category = ''
        for row in current_df.iloc[4:].itertuples(index=False):
            if row[0] == "RANK" or pd.isna(row[0]) or (isinstance(row[0], int) and pd.isna(row[1])) or (isinstance(row[1], str) and row[1].isspace()):
                continue
            if isinstance(row[0], str):
                current_weight_category = row[0]
                tournaments[i]['weight_categories'][current_weight_category] = []
                continue
            athlete = { key : row[i] for i, key in enumerate(headlines) }
            tournaments[i]['weight_categories'][current_weight_category].append(athlete)

    to_json[name] = tournaments
    
with open('Соревнования.json', 'w', encoding="utf-8") as f:
    json.dump(to_json, f, indent=4)


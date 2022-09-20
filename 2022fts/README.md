
1. put the extracted 2022fts/datasets/urbandict-word-defs.csv
   download the file from https://www.kaggle.com/datasets/therohk/urban-dictionary-words-dataset


total records: 2_580_925

TODOs:
[ ] BulkInsert per 2000 (first 1m), 10000 (next 1m), the rest (reindex once at the end only)
    duration per record
[ ] Search first word linearly
[ ] Search last word of each dataset linearly
[ ] Search duration by locale
[ ] Deleting first 100 records
[ ] Check disk usage and memory usage of the docker-compose
[ ] Create makefile to clean docker data files
[ ] Add gitignore for data files

import sqlite3
import pandas as pd
import os

def create_database_from_csv(csv_file_path):
    """
    Creates a SQLite database from a CSV file.
    The table name will be derived from the CSV filename.
    """
    # Read the CSV file
    df = pd.read_csv(csv_file_path)
    
    # Get the base filename without extension to use as table name
    table_name = os.path.splitext(os.path.basename(csv_file_path))[0].upper()
    
    # Create database connection
    db_name = f"data/{table_name}.db"
    conn = sqlite3.connect(db_name)
    
    # Write the dataframe to SQLite
    df.to_sql(table_name, conn, if_exists='replace', index=False)
    
    # Get column names and types for the prompt
    columns_info = []
    for column in df.columns:
        dtype = str(df[column].dtype)
        if dtype.startswith('int'):
            sql_type = 'INTEGER'
        elif dtype.startswith('float'):
            sql_type = 'FLOAT'
        else:
            sql_type = 'TEXT'
        columns_info.append(f"{column} ({sql_type})")
    
    conn.close()
    
    return table_name, db_name, columns_info

def read_sql_query(sql, db):
    """
    Execute a SQL query and return the results
    """
    conn = sqlite3.connect(db)
    cur = conn.cursor()
    cur.execute(sql)
    rows = cur.fetchall()
    conn.close()
    return rows

def get_table_info(db_path, table_name):
    """
    Get table schema information
    """
    conn = sqlite3.connect(db_path)
    cur = conn.cursor()
    cur.execute(f"PRAGMA table_info({table_name})")
    columns = cur.fetchall()
    conn.close()
    return columns
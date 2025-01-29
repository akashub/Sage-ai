from dotenv import load_dotenv
import streamlit as st
import os
import google.generativeai as genai
import pandas as pd
import sqlite3
from tempfile import NamedTemporaryFile

load_dotenv()

genai.configure(api_key=os.getenv("GOOGLE_API_KEY"))

def get_gemini_response(question, prompt):
    model = genai.GenerativeModel('gemini-pro')
    response = model.generate_content([prompt, question])
    return response.text

def create_database_from_dataframe(df, table_name):
    """
    Creates a SQLite database from a pandas DataFrame
    """
    # Create temporary database in memory
    conn = sqlite3.connect(':memory:')
    
    # Write the dataframe to SQLite
    df.to_sql(table_name, conn, if_exists='replace', index=False)
    
    # Get column information
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
    
    return conn, columns_info

def read_sql_query(sql, conn):
    """
    Execute a SQL query and return the results
    """
    cur = conn.cursor()
    cur.execute(sql)
    return cur.fetchall()

def create_prompt(table_name, columns_info):
    return f"""
    You are an expert in converting English questions to SQL query!
    The SQL database has the table name {table_name} with the following columns:
    {', '.join(columns_info)}
    
    For example:
    Example 1 - How many entries are present?
    The SQL command will be: SELECT COUNT(*) FROM {table_name}
    
    Example 2 - What are the top 5 movies by vote average?
    The SQL command will be: SELECT title, vote_average FROM {table_name} ORDER BY vote_average DESC LIMIT 5
    
    The sql code should not have ``` in beginning or end and sql word in output.
    Return only the SQL query without any additional text or explanation.
    """

# Streamlit App
st.set_page_config(page_title="Sage-AI V1 Demo: Dynamic SQL Query Assistant")
st.title("Sage-AI V1: LLM-Based Dynamic SQL Query Assistant")

# File uploader
uploaded_file = st.file_uploader("Choose a CSV file", type="csv")

if uploaded_file is not None:
    # Read the CSV directly into a DataFrame
    df = pd.read_csv(uploaded_file)
    
    # Get table name from file name without extension
    table_name = os.path.splitext(uploaded_file.name)[0].upper()
    
    # Create in-memory database and get column info
    conn, columns_info = create_database_from_dataframe(df, table_name)
    
    # Create prompt with table information
    prompt = create_prompt(table_name, columns_info)
    
    # User input
    question = st.text_input("What would you like to know about the data?")
    submit = st.button("Ask Question")
    
    if submit and question:
        try:
            # Get SQL query from Gemini
            sql_query = get_gemini_response(question, prompt)
            
            # Display the SQL query
            st.code(sql_query, language="sql")
            
            # Execute query and display results
            response = read_sql_query(sql_query, conn)
            
            # Display results in a dataframe
            if response:
                df_response = pd.DataFrame(response, columns=df.columns[:len(response[0])])
                st.dataframe(df_response)
            else:
                st.write("No results found.")
                
        except Exception as e:
            st.error(f"An error occurred: {str(e)}")
else:
    st.info("Please upload a CSV file to begin.")
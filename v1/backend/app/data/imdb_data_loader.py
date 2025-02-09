# app/data/imdb_data_loader.py
import pandas as pd
from sqlalchemy import create_engine, text
from typing import Optional
import json
import numpy as np

class IMDBDataLoader:
    def __init__(self, database_url: str):
        self.database_url = database_url
        self.engine = create_engine(database_url)

    def load_csv_to_df(self, file_path: str) -> pd.DataFrame:
        """Load IMDB CSV file into a pandas DataFrame"""
        df = pd.read_csv(file_path)
        
        # Clean and preprocess the data
        df = self._preprocess_dataframe(df)
        
        return df

    def _preprocess_dataframe(self, df: pd.DataFrame) -> pd.DataFrame:
        """Preprocess the dataframe"""
        # Handle date conversion
        df['release_date'] = pd.to_datetime(df['release_date'], errors='coerce').dt.date
        
        # Convert boolean values
        if 'adult' in df.columns:
            df['adult'] = df['adult'].map({'True': True, 'False': False})
        
        # Handle numeric columns
        numeric_columns = ['vote_average', 'vote_count', 'revenue', 'budget', 'runtime', 'popularity']
        for col in numeric_columns:
            if col in df.columns:
                df[col] = pd.to_numeric(df[col], errors='coerce')
                df[col] = df[col].fillna(0)
        
        # Handle JSON-like string columns
        json_columns = ['genres', 'production_companies', 'production_countries', 
                       'spoken_languages', 'keywords']
        
        for col in json_columns:
            if col in df.columns:
                df[col] = df[col].apply(self._process_json_field)
        
        return df

    def _process_json_field(self, field):
        """Process JSON-like fields, ensuring they're valid JSON strings"""
        if pd.isna(field) or field == '':
            return '[]'
        
        try:
            # If it's already a string representation of JSON, return it
            if isinstance(field, str):
                json.loads(field)  # Validate JSON
                return field
            # If it's a Python object, convert to JSON string
            return json.dumps(field)
        except:
            return '[]'

    def create_database_schema(self):
        """Create the movies table in the database"""
        create_table_sql = """
        CREATE TABLE IF NOT EXISTS movies (
            id INTEGER PRIMARY KEY,
            title TEXT,
            vote_average FLOAT,
            vote_count INTEGER,
            status TEXT,
            release_date DATE,
            revenue BIGINT,
            runtime INTEGER,
            adult BOOLEAN,
            budget BIGINT,
            imdb_id TEXT,
            original_language TEXT,
            original_title TEXT,
            overview TEXT,
            popularity FLOAT,
            tagline TEXT,
            genres TEXT,
            production_companies TEXT,
            production_countries TEXT,
            spoken_languages TEXT,
            keywords TEXT
        );
        """
        with self.engine.connect() as conn:
            conn.execute(text(create_table_sql))
            conn.commit()

    def load_data_to_db(self, df: pd.DataFrame, if_exists: str = 'replace'):
        """Load DataFrame into the database"""
        # Ensure all text columns are string type
        text_columns = df.select_dtypes(include=['object']).columns
        for col in text_columns:
            df[col] = df[col].astype(str)
        
        # Load to database
        df.to_sql('movies', self.engine, if_exists=if_exists, index=False)
        
        # Create indexes for common query columns
        with self.engine.connect() as conn:
            conn.execute(text("CREATE INDEX IF NOT EXISTS idx_release_date ON movies(release_date)"))
            conn.execute(text("CREATE INDEX IF NOT EXISTS idx_vote_average ON movies(vote_average)"))
            conn.execute(text("CREATE INDEX IF NOT EXISTS idx_popularity ON movies(popularity)"))
            conn.commit()

    def validate_data(self, df: pd.DataFrame) -> dict:
        """Validate the loaded data"""
        validation_results = {
            "total_rows": len(df),
            "null_counts": df.isnull().sum().to_dict(),
            "data_types": df.dtypes.astype(str).to_dict(),
            "sample_rows": df.head(5).to_dict('records'),
            "issues": []
        }
        
        # Check for required columns
        required_columns = ['id', 'title', 'vote_average', 'release_date']
        missing_columns = [col for col in required_columns if col not in df.columns]
        if missing_columns:
            validation_results["issues"].append(f"Missing required columns: {missing_columns}")
        
        # Check for data quality issues
        if 'vote_average' in df.columns:
            invalid_ratings = df[~df['vote_average'].between(0, 10)].shape[0]
            if invalid_ratings > 0:
                validation_results["issues"].append(f"Found {invalid_ratings} ratings outside valid range (0-10)")
        
        if 'release_date' in df.columns:
            null_dates = df['release_date'].isnull().sum()
            if null_dates > 0:
                validation_results["issues"].append(f"Found {null_dates} null release dates")
        
        return validation_results
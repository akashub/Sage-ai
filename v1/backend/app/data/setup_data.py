# app/data/setup_data.py
import os
from pathlib import Path
import sys
from sqlalchemy import create_engine
from imdb_data_loader import IMDBDataLoader

def setup_database():
    # Get the project root directory
    current_dir = Path(__file__).parent
    project_root = current_dir.parent.parent  # This should point to the project root
    
    # Set up database URL - using SQLite for simplicity
    database_path = project_root / "data" / "imdb.db"
    database_url = f"sqlite:///{database_path}"
    
    # Initialize loader
    loader = IMDBDataLoader(database_url=database_url)
    
    try:
        # Create schema
        loader.create_database_schema()
        print("Schema created successfully")
        
        # Path to the CSV file
        csv_path = current_dir / "Imdb Movie Dataset.csv"
        
        if not csv_path.exists():
            raise FileNotFoundError(f"CSV file not found at {csv_path}")
        
        # Load and validate data
        print("Loading CSV file...")
        df = loader.load_csv_to_df(str(csv_path))
        
        # Validate data
        print("Validating data...")
        validation_results = loader.validate_data(df)
        print("\nValidation Results:")
        print(f"Total rows: {validation_results['total_rows']}")
        print("\nNull counts:")
        for col, count in validation_results['null_counts'].items():
            if count > 0:
                print(f"{col}: {count}")
        
        if validation_results['issues']:
            print("\nIssues found:")
            for issue in validation_results['issues']:
                print(f"- {issue}")
        
        # Load to database
        print("\nLoading data to database...")
        loader.load_data_to_db(df)
        print("Data loaded successfully!")
        
        # Test query
        print("\nTesting database connection...")
        engine = create_engine(database_url)
        with engine.connect() as conn:
            result = conn.execute("SELECT COUNT(*) FROM movies").scalar()
            print(f"Total records in database: {result}")
        
        return True
        
    except Exception as e:
        print(f"Error: {str(e)}")
        return False

if __name__ == "__main__":
    success = setup_database()
    sys.exit(0 if success else 1)
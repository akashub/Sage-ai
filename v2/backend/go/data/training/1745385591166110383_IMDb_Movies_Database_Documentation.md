# IMDb Movies Database Documentation

This document provides comprehensive information about the IMDb movies database schema, its tables, relationships, and recommended query patterns.

## Overview

The IMDb movies database contains information about movies including their basic details, ratings, financial data, cast, crew, and categorization. The database is designed with a normalized structure to efficiently store and query movie data.

## Tables and Their Relationships

### Main Tables

#### `movies`
The central table containing core movie information.
- `id`: Primary key, unique identifier for each movie.
- `title`: The title of the movie as commonly known.
- `original_title`: The title in its original language.
- `original_language`: Code representing the movie's original language (e.g., "en" for English).
- `overview`: Plot summary or description.
- `tagline`: Short promotional phrase associated with the movie.
- `status`: Current state of the movie (e.g., "Released", "In Production").
- `release_date`: Official release date.
- `runtime`: Duration in minutes.
- `budget`: Production cost (USD).
- `revenue`: Total earnings (USD).
- `popularity`: Metric indicating the movie's popularity.
- `vote_average`: Average rating (scale 0-10).
- `vote_count`: Number of ratings received.
- `adult`: Boolean indicating if the movie contains adult content.
- `imdb_id`: IMDb's unique identifier for the movie.

### Lookup Tables and Many-to-Many Relationships

#### `genres` and `movie_genres`
Movie categorization by genre.
- A movie can have multiple genres.
- A genre can be associated with multiple movies.

#### `production_companies` and `movie_production_companies`
Companies involved in producing the movies.
- A movie can be produced by multiple companies.
- A company can produce multiple movies.

#### `countries` and `movie_production_countries`
Countries where movies were produced.
- A movie can be produced in multiple countries.
- A country can be involved in the production of multiple movies.

#### `languages` and `movie_spoken_languages`
Languages spoken in the movies.
- A movie can feature multiple languages.
- A language can be used in multiple movies.

#### `keywords` and `movie_keywords`
Terms or phrases associated with the movies for categorization.
- A movie can have multiple keywords.
- A keyword can be associated with multiple movies.

## Common Query Patterns

### Basic Movie Retrieval
```sql
SELECT * FROM movies WHERE title LIKE '%Star Wars%';
```

### Finding Movies by Genre
```sql
SELECT m.* 
FROM movies m
JOIN movie_genres mg ON m.id = mg.movie_id
JOIN genres g ON mg.genre_id = g.genre_id
WHERE g.genre_name = 'Action';
```

### Top-Rated Movies
```sql
SELECT title, vote_average, vote_count
FROM movies
WHERE vote_count > 1000
ORDER BY vote_average DESC
LIMIT 10;
```

### Highest-Grossing Movies
```sql
SELECT title, revenue, budget, (revenue - budget) AS profit
FROM movies
WHERE budget > 0 AND revenue > 0
ORDER BY revenue DESC
LIMIT 10;
```

### Movies by Production Company
```sql
SELECT m.title, m.release_date
FROM movies m
JOIN movie_production_companies mpc ON m.id = mpc.movie_id
JOIN production_companies pc ON mpc.company_id = pc.company_id
WHERE pc.company_name = 'Warner Bros.';
```

### Movies by Language
```sql
SELECT m.title
FROM movies m
JOIN movie_spoken_languages msl ON m.id = msl.movie_id
JOIN languages l ON msl.language_id = l.language_id
WHERE l.language_name = 'French';
```

### Movies by Production Country
```sql
SELECT m.title
FROM movies m
JOIN movie_production_countries mpc ON m.id = mpc.movie_id
JOIN countries c ON mpc.country_id = c.country_id
WHERE c.country_name = 'Japan';
```

### Movies with Specific Keywords
```sql
SELECT m.title
FROM movies m
JOIN movie_keywords mk ON m.id = mk.movie_id
JOIN keywords k ON mk.keyword_id = k.keyword_id
WHERE k.keyword_name IN ('superhero', 'based on comic');
```

### Movie Trends by Year
```sql
SELECT EXTRACT(YEAR FROM release_date) AS year, 
       COUNT(*) AS movie_count, 
       AVG(vote_average) AS avg_rating
FROM movies
GROUP BY year
ORDER BY year DESC;
```

### Most Popular Genres Over Time
```sql
SELECT EXTRACT(YEAR FROM m.release_date) AS year, 
       g.genre_name, 
       COUNT(*) AS movie_count
FROM movies m
JOIN movie_genres mg ON m.id = mg.movie_id
JOIN genres g ON mg.genre_id = g.genre_id
WHERE m.release_date IS NOT NULL
GROUP BY year, g.genre_name
ORDER BY year DESC, movie_count DESC;
```

## Performance Considerations

1. The schema includes indexes on frequently queried columns like `release_date`, `vote_average`, `popularity`, `revenue`, and `budget`.
2. For text searches on `title` or `overview`, consider using full-text search capabilities of your database system.
3. When querying across multiple many-to-many relationships, use appropriate JOIN strategies to optimize performance.
4. For analytical queries on large datasets, consider materialized views or pre-aggregated tables.

## Data Integrity

1. The schema enforces referential integrity through foreign key constraints.
2. The `movies` table uses `id` as its primary key, while the `imdb_id` field provides a unique external reference.
3. Lookup tables use surrogate keys with meaningful name fields to support internationalization and avoid repetition.

## Recommended Extensions

1. Add a `cast` and `crew` structure to track actors, directors, and other film personnel.
2. Implement a `collections` table to group movie franchises (e.g., "Star Wars Saga").
3. Add user-specific tables for personalized recommendations and watch history.

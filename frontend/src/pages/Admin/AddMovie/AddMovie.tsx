import React, {useEffect, useState} from "react";
import '@assets/styles/pages/admin/addmovie.css';
import FileDropZone from "@components/specific/FileDropZone.tsx";
import {addMovies, getSearchMoviesAPI, Movies} from "@components/gRPC.tsx";

// Компонент Home
const AddMovie: React.FC = () => {
    const [query, setQuery] = useState('');
    const [movies, setMovies] = useState<Movies[]>([]);
    const [selectedMovie, setSelectedMovie] = useState<Movies | null>(null);
    const [isFileUploaded, setIsFileUploaded] = useState<boolean>(false);
    const [debounceTimeout, setDebounceTimeout] = useState<ReturnType<typeof setTimeout> | null>(null);

    useEffect(() => {
        if (debounceTimeout) {
            clearTimeout(debounceTimeout);
        }

        // Устанавливаем таймер на 2 секунды
        const timeoutId = setTimeout(async () => {
            if (query) {
                try {
                    const result = await fetchMovies(query);
                    console.log('Fetched movies:', result); // Проверка данных
                    if (result != undefined) {
                        setMovies(result);
                    }
                } catch (error) {
                    console.error('Error fetching movies:', error);
                }
            }
        }, 2000);

        setDebounceTimeout(timeoutId);

        // Очистка таймера при размонтировании компонента или изменении query
        return () => clearTimeout(timeoutId);
    }, [query]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setQuery(e.target.value);
    };

    const handleMovieSelect = (movie: Movies) => {
        setSelectedMovie(movie);
    };

    const handleFileUploaded = () => {
        setIsFileUploaded(true);
    };

    const fetchMovies = async (query: string): Promise<Movies[] | undefined> => {
        let data;
        try {
            data = await getSearchMoviesAPI(5, query);
        } catch (error) {
            console.log(error)
        }

        return data;
    };

    const handleSubmit = async () => {
        if (selectedMovie && isFileUploaded) {
            const {title, description, typeMovie, releaseDate, scoreKP, genres} = selectedMovie;

            try {
                await addMovies(title, description, releaseDate, scoreKP, typeMovie, genres);
            } catch (error) {
                console.log(error)
            }
        }
    };

    return (
        <div className="sections">
            <div className="section__movie-add">
                <input type="text" list="movies" placeholder="Название фильма" id="section__movies-input" onChange={handleChange}/>
                {movies.length > 0 && (
                    <div className="section__movie-results">
                        <ul>
                            {movies.map((movie) => (
                                <li
                                    key={movie.id}
                                    onClick={() => handleMovieSelect(movie)}
                                    className={selectedMovie?.id === movie.id ? 'active' : ''}
                                >{movie.title}</li>
                            ))}
                        </ul>
                    </div>
                )}
                <FileDropZone onFileUploaded={handleFileUploaded} />
                <button onClick={handleSubmit} disabled={!selectedMovie || !isFileUploaded}>
                    <div id="circle"></div>
                    <span>Отправить</span>
                </button>
            </div>
        </div>
    );
};

export default AddMovie;

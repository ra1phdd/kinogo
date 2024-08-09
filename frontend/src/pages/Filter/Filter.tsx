import React, { useEffect } from "react";
import SearchAside from "@components/specific/Search.tsx";
import FilterAside from "@components/specific/Filter.tsx";
import BestMovieAside from "@components/specific/BestMovie.tsx";
import { useLocation, useSearchParams } from "react-router-dom";
import MovieCard from "@components/common/MovieCard.tsx";
import '@assets/styles/pages/filter.css';
import useFilterMovies from '@/hooks/fetchMovies/useFilterMovies.ts';
import useScroll from '@/hooks/useScroll';

const Filter: React.FC = () => {
    const location = useLocation();
    const formData = location.state as { genres: string[], year__min: string, year__max: string };
    const [searchParams] = useSearchParams();
    const text = searchParams.get('text') || "";

    const genres = formData?.genres?.join(', ') || '';
    const yearMin = Number(formData?.year__min || 0);
    const yearMax = Number(formData?.year__max || 9999);

    const { movies, loading, loadMovies } = useFilterMovies(1, genres, yearMin, yearMax, text, location.pathname);

    useScroll(loadMovies);

    useEffect(() => {
        loadMovies();
    }, []);

    return (
        <>
            <div className="sections">
                <div className="section__filters">
                    {text === "" ? (
                        <>
                            <h3>Фильтры: </h3>
                            {genres && <span>{genres}</span>}
                            {yearMax && <span>{"<"}{yearMax}</span>}
                            {yearMin && <span>{">"}{yearMin}</span>}
                        </>
                    ) : (
                        <>
                            <h3>Результаты поиска по запросу: </h3>
                            <span>{text}</span>
                        </>
                    )}
                </div>
                <div className="section__videos">
                    {movies.map((movie, index) => (
                        <MovieCard key={movie.id} movie={movie} index={index} limit={12} />
                    ))}
                    {loading && <p>Загрузка...</p>}
                    {movies.length === 0 && !loading && <p>Данных в БД нет</p>}
                </div>
            </div>
            <aside>
                <SearchAside />
                <FilterAside />
                <BestMovieAside />
            </aside>
        </>
    );
}

export default Filter;

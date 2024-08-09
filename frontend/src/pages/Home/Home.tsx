import React, { useEffect } from "react";
import MovieCard from "@components/common/MovieCard.tsx";
import '@assets/styles/pages/home.css'
import useScroll from "@/hooks/useScroll.ts";
import useMovies from "@/hooks/fetchMovies/useMovies.ts";

// Компонент Home
const Home: React.FC = () => {
    const { movies, loading, loadMovies } = useMovies(1);

    useScroll(loadMovies);

    useEffect(() => {
        loadMovies();
    }, []);

    return (
        <>
            <div className="sections">
                <div className="section__videos">
                    {movies.map((movie, index) => (
                        <MovieCard key={movie.id} movie={movie} index={index} limit={15}  />
                    ))}
                </div>
            </div>
            {loading && <div>Загрузка...</div>}
        </>
    );
};

export default Home;

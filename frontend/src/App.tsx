import { BrowserRouter, Routes, Route } from "react-router-dom";
import '@assets/styles/vendor/animate.min.css'
import Navigation from "@components/common/Navigation.tsx";
import Auth from "@components/common/TelegramAuth.tsx";
import {lazy, useEffect} from "react";
import Cookies from "js-cookie";
import { v4 as uuidv4 } from 'uuid';
import {metricNewUser} from "@components/gRPC.tsx";
import TimeTracker from "@components/common/TimeTracker.tsx";

// Ленивая загрузка страниц
const Home = lazy(() => import("./pages/Home/Home.tsx"));
const Movie = lazy(() => import("./pages/Movie/Movie.tsx"));
const Filter = lazy(() => import("@/pages/Filter/Filter.tsx"));
const Movies = lazy(() => import("@/pages/Movies/Movies.tsx"));
const AddMovie = lazy(() => import("@/pages/Admin/AddMovie/AddMovie.tsx"));

function App() {
    useEffect(() => {
        let userUUID = Cookies.get('userUUID');

        if (userUUID == undefined) {
            userUUID = uuidv4();

            try {
                metricNewUser();
            } catch (error) {
                console.error('Error fetching movie:', error);
            }
        }

        Cookies.set('userUUID', userUUID, { expires: 365 });
    }, []);

    return (
    <>
        <TimeTracker/>
        <BrowserRouter>
            <div className="container">
                <header>
                    <div className="header__logotype">
                        <h1><a href="/">KINOTEATR</a></h1>
                    </div>
                    <nav>
                        <Navigation/>
                    </nav>
                    <div className="header__auth">
                        <Auth/>
                    </div>
                </header>
                <main>
                    <Routes>
                        {/* Мэйн */}
                        <Route path="/" element={<Home />} />
                        <Route path="/films" element={<Movies />} />
                        <Route path="/cartoons" element={<Movies />} />
                        <Route path="/telecasts" element={<Movies />} />
                        <Route path="/id/:id" element={<Movie />} />
                        <Route path="/search" element={<Filter />} />
                        <Route path="/filter" element={<Filter />} />
                        {/* Админ-панель */}
                        <Route path="/admin/addmovie" element={<AddMovie />} />
                    </Routes>
                </main>
            </div>
        </BrowserRouter>
    </>
  )
}

export default App

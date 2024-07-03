import { BrowserRouter, Routes, Route } from "react-router-dom";
import '@assets/styles/vendor/animate.min.css'
import Navigation from "@components/Navigation.tsx";
import Auth from "@components/TelegramAuth.tsx";
import {lazy} from "react";
import Studio from "@/pages/Studio/Studio.tsx";

// Ленивая загрузка страниц
const Home = lazy(() => import("./pages/Home/Home.tsx"));
const Movie = lazy(() => import("./pages/Movie/Movie.tsx"));
const Filter = lazy(() => import("@/pages/Filter/Filter.tsx"));
const Movies = lazy(() => import("@/pages/Movies/Movies.tsx"));

function App() {
    return (
    <>
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
                        <Route path="/studio" element={<Studio />} />
                    </Routes>
                </main>
            </div>
        </BrowserRouter>
    </>
  )
}

export default App

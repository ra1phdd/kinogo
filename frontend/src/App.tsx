import '@components/Navigation.tsx'
import {BrowserRouter, Routes, Route} from "react-router-dom";
import '@assets/css/src/main.css'
import '@assets/css/dist/animate.min.css'
import Navigation from "@components/Navigation.tsx";
import Home from "./pages/Home/Home.tsx"
import Movie from "./pages/Movie/Movie.tsx"
import Auth from "@components/TelegramAuth.tsx"
import Search from "@/pages/Search/Search.tsx";
import Filter from "@/pages/Filter/Filter.tsx";

function App() {
    return (
    <>
        <BrowserRouter>
            <div className="container">
                <header className="animate__animated animate__fadeInDown animate__faster" style={{overflow: 'hidden'}}>
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
                        <Route path="/" element={<Home />} />
                        <Route path="/films" element={<Home />} />
                        <Route path="/cartoons" element={<Home />} />
                        <Route path="/telecasts" element={<Home />} />
                        <Route path="/id/:id" element={<Movie />} />
                        <Route path="/search" element={<Search />} />
                        <Route path="/filter" element={<Filter />} />
                    </Routes>
                </main>
            </div>
        </BrowserRouter>
    </>
  )
}

export default App

import Home from "./screen/Home";
import Test from "./screen/Test";
import { Routes, Route } from "react-router-dom";

function Router() {
    return (
        <Routes>
            <Route path="/" element={<Home />} />
            <Route path="test" element={<Test />} />
      </Routes>

    );
}

export default Router;
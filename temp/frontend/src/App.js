import './style_app.css';
import Activities from './activities';
import Login from './login';
import Activity from './activity';
import MyActivities from './MyActivities';
import Inscription from './inscription';
import MyActivity from './myactivity';
import Admin from './admin';
import ActivityAdmin from './activityadmin';
import CreateActivity from "./createactivity";
import EditActivity from './editactivity';
import { BrowserRouter, Routes, Route } from "react-router-dom";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/login" element={<Login />} />
        <Route path="/activities" element={<Activities />} />
        <Route path="/activities/:id" element={<Activity />} />
        <Route path="/users/activities" element={<MyActivities />} />
        <Route path="/users/inscription/:id" element={<Inscription />} />
        <Route path="/myactivities/:id" element={<MyActivity />} />
        <Route path="/users/admin" element={<Admin />} />
        <Route path="/createactivity" element={<CreateActivity />} />
        <Route path="/activityadmin/:id" element={<ActivityAdmin />} />
        <Route path="/editactivity/:id" element={<EditActivity />} />

      </Routes>
    </BrowserRouter>
  );  
}

export default App;

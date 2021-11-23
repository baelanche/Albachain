import passport from "passport"
import Worker from "./model/worker"

passport.use(Worker.createStrategy());
passport.serializeUser(Worker.serializeUser());
passport.deserializeUser(Worker.deserializeUser());
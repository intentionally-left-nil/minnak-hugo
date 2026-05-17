// MiNNaK Hugo Theme — entry point
// Imports sidebar behaviour; CSS is handled separately via Hugo's asset pipeline.

import sidebar from "./sidebar.js";
import fixSamePageSearch from "./search.js";

sidebar();
fixSamePageSearch();

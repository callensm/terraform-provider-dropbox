package dropbox

const fileIDPattern = "((/|id:).*|nspath:[0-9]+:.*)|ns:[0-9]+(/.*)?"

const emailPattern = "^['&A-Za-z0-9._%+-]+@[A-Za-z0-9-][A-Za-z0-9.-]*.[A-Za-z]{2,15}$"

const folderPathPattern = "(/(.|[\r\n])*)|(ns:[0-9]+(/.*)?)"

const uploadPathPattern = "(/(.|[\r\n])*)|(ns:[0-9]+(/.*)?)|(id:.*)"

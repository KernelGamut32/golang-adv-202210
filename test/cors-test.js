function main() 
{
    console.log("call that require cors enabled");
    $.ajax({
        url: "http://localhost:5000/user",
        success: function(data) 
        {
            console.log("log response on success");
            console.log(data);
        },
        error: function(err)
        {
            console.log("log message on failure");
            console.log(err);
        }
    });
}
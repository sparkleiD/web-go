function disableEnterKey(event) {
    if (event.key === 'Enter') {
        event.preventDefault();
    }
}
document.addEventListener('keypress', disableEnterKey);

function signupChange1() {
    var div = document.getElementById("signup");
    div.style.width = "600px";
    var che = document.getElementById("signupcheck");
    var ins = document.getElementById("signupins");
    che.style.display = "block";
    ins.style.display = "none";
    loginChange2();
}
function signupChange2() {
    var div = document.getElementById("signup");
    div.style.width = "300px";
    var che = document.getElementById("signupcheck");
    var ins = document.getElementById("signupins");
    che.style.display = "none";
    ins.style.display = "block";
}
function loginChange1() {
    var div = document.getElementById("login");
    div.style.width = "600px";
    var che = document.getElementById("logincheck");
    var ins = document.getElementById("loginins");
    che.style.display = "block";
    ins.style.display = "none";
    signupChange2();
}
function loginChange2() {
    var div = document.getElementById("login");
    div.style.width = "300px";
    var che = document.getElementById("logincheck");
    var ins = document.getElementById("loginins");
    che.style.display = "none";
    ins.style.display = "block";
}
//以上为动态网页代码
function Ajax(url, content) {
    var request = new XMLHttpRequest();
    request.open('POST', url, false); //false表示同步请求
    request.send(content);
    return request.responseText;
}

S_email.onchange = function () {
    var S_email = this.value;
    var reg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/;
    if (!reg.test(S_email)) {
        alert("邮箱格式不正确，请重新输入！");
        document.getElementById("S_email").value = "";
        return false;
    }
    var formdata = new FormData;
    formdata.append("email", S_email);
    var ajax = Ajax("/ajax/userinfo", formdata);
    if (ajax == "false") {
        alert("邮箱已被注册，请登陆！");
        document.getElementById("S_email").value = "";
        return false;
    }
};
L_email.onchange = function () {
    var L_email = this.value;
    var reg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/;
    if (!reg.test(L_email)) {
        alert("邮箱格式不正确，请重新输入！");
        document.getElementById("L_email").value = "";
        return false;
    }
};
S_user.onchange = function () {
    var S_user = this.value;
    var reg = /^\S{1,10}$/;
    if (!reg.test(S_user)) {
        alert("用户名长度不能超过10位！");
        document.getElementById("S_user").value = "";
        return false;
    }
};
S_pwd.onchange = function () {
    var S_pwd = this.value;
    var reg = /^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,}$/;
    if (!reg.test(S_pwd)) {
        alert("密码长度至少为6位，且要由数字和字母组成,请重新输入！");
        document.getElementById("S_pwd").value = "";
        return false;
    }
};
function checkpassword() {
    var C_pwd = document.getElementById("S_pwd").value;
    var C_repwd = document.getElementById("S_repwd").value;
    if (C_pwd == C_repwd) {
        document.getElementById("tishi").innerHTML =
            "<font color='green'>两次密码输入一致</font>";
    } else {
        document.getElementById("tishi").innerHTML =
            "<font color='red'>两次输入密码不一致!</font>";
    }
}
//以上是检测内容是否符合格式的代码

function signupscuueed() {
    localStorage.setItem("username", document.getElementById("S_user").value);
    localStorage.setItem("password", document.getElementById("S_pwd").value);
    localStorage.setItem("email", document.getElementById("S_email").value);

    var S_Email = document.getElementById("S_email").value;
    var S_User = document.getElementById("S_user").value;
    var S_Iconname = document.getElementById("toux").value.split("\\").pop();
    var S_Pwd = document.getElementById("S_pwd").value;
    var S_Repwd = document.getElementById("S_repwd").value;
    if (S_Email == "" || S_User == "" || S_Iconname == "" || S_Pwd == "") {
        alert("邮箱、用户名、头像或密码不能为空！");
        return false;
    } else if (S_Pwd != S_Repwd) {
        alert("请输入一致的密码！");
        return false;
    } 
}

function loginsucceed() {
    var L_Email = document.getElementById("L_email").value;
    var L_Pwd = document.getElementById("L_pwd").value;
    if (L_Email == "" || L_Pwd == "") {
        alert("邮箱或密码不能为空！");
        return false;
    }
}
//以上为表单提交校验代码

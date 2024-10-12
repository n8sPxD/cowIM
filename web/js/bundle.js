(function(){function r(e,n,t){function o(i,f){if(!n[i]){if(!e[i]){var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a}var p=n[i]={exports:{}};e[i][0].call(p.exports,function(r){var n=e[i][1][r];return o(n||r)},p,p.exports,r,e,n,t)}return n[i].exports}for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o}return r})()({1:[function(require,module,exports){
const cowsay = require("cowsay");

const cowText = cowsay.say({
    text: "Welcome to Cow IM!\n\nCow is not Copy-On-Write!",
    e: "^^",
    T: "U "
});

document.getElementById("cowsay-container").innerText = cowText;
},{"cowsay":2}],2:[function(require,module,exports){
(function (global, factory) {
	typeof exports === 'object' && typeof module !== 'undefined' ? factory(exports) :
	typeof define === 'function' && define.amd ? define(['exports'], factory) :
	(global = typeof globalThis !== 'undefined' ? globalThis : global || self, factory(global.cowsay = {}));
}(this, (function (exports) { 'use strict';

	var ansiRegex = () => {
		const pattern = [
			'[\\u001B\\u009B][[\\]()#;?]*(?:(?:(?:(?:;[-a-zA-Z\\d\\/#&.:=?%@~_]+)*|[a-zA-Z\\d]+(?:;[a-zA-Z\\d]*)*)?\\u0007)',
			'(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))'
		].join('|');

		return new RegExp(pattern, 'g');
	};

	var stripAnsi = input => typeof input === 'string' ? input.replace(ansiRegex(), '') : input;

	/* eslint-disable yoda */
	var isFullwidthCodePoint = x => {
		if (Number.isNaN(x)) {
			return false;
		}

		// code points are derived from:
		// http://www.unix.org/Public/UNIDATA/EastAsianWidth.txt
		if (
			x >= 0x1100 && (
				x <= 0x115f ||  // Hangul Jamo
				x === 0x2329 || // LEFT-POINTING ANGLE BRACKET
				x === 0x232a || // RIGHT-POINTING ANGLE BRACKET
				// CJK Radicals Supplement .. Enclosed CJK Letters and Months
				(0x2e80 <= x && x <= 0x3247 && x !== 0x303f) ||
				// Enclosed CJK Letters and Months .. CJK Unified Ideographs Extension A
				(0x3250 <= x && x <= 0x4dbf) ||
				// CJK Unified Ideographs .. Yi Radicals
				(0x4e00 <= x && x <= 0xa4c6) ||
				// Hangul Jamo Extended-A
				(0xa960 <= x && x <= 0xa97c) ||
				// Hangul Syllables
				(0xac00 <= x && x <= 0xd7a3) ||
				// CJK Compatibility Ideographs
				(0xf900 <= x && x <= 0xfaff) ||
				// Vertical Forms
				(0xfe10 <= x && x <= 0xfe19) ||
				// CJK Compatibility Forms .. Small Form Variants
				(0xfe30 <= x && x <= 0xfe6b) ||
				// Halfwidth and Fullwidth Forms
				(0xff01 <= x && x <= 0xff60) ||
				(0xffe0 <= x && x <= 0xffe6) ||
				// Kana Supplement
				(0x1b000 <= x && x <= 0x1b001) ||
				// Enclosed Ideographic Supplement
				(0x1f200 <= x && x <= 0x1f251) ||
				// CJK Unified Ideographs Extension B .. Tertiary Ideographic Plane
				(0x20000 <= x && x <= 0x3fffd)
			)
		) {
			return true;
		}

		return false;
	};

	var stringWidth = str => {
		if (typeof str !== 'string' || str.length === 0) {
			return 0;
		}

		str = stripAnsi(str);

		let width = 0;

		for (let i = 0; i < str.length; i++) {
			const code = str.codePointAt(i);

			// Ignore control characters
			if (code <= 0x1F || (code >= 0x7F && code <= 0x9F)) {
				continue;
			}

			// Ignore combining characters
			if (code >= 0x300 && code <= 0x36F) {
				continue;
			}

			// Surrogates
			if (code > 0xFFFF) {
				i++;
			}

			width += isFullwidthCodePoint(code) ? 2 : 1;
		}

		return width;
	};

	var say = function (text, wrap) {
		var delimiters = {
			first : ["/", "\\"],
			middle : ["|", "|"],
			last : ["\\", "/"],
			only : ["<", ">"]
		};

		return format(text, wrap, delimiters);
	};

	var think = function (text, wrap) {
		var delimiters = {
			first : ["(", ")"],
			middle : ["(", ")"],
			last : ["(", ")"],
			only : ["(", ")"]
		};

		return format(text, wrap, delimiters);
	};

	function format (text, wrap, delimiters) {
		var lines = split(text, wrap);
		var maxLength = max(lines);

		var balloon;
		if (lines.length === 1) {
			balloon = [
				" " + top(maxLength),
				delimiters.only[0] + " " + lines[0] + " " + delimiters.only[1],
				" " + bottom(maxLength)
			];
		} else {
			balloon = [" " + top(maxLength)];

			for (var i = 0, len = lines.length; i < len; i += 1) {
				var delimiter;

				if (i === 0) {
					delimiter = delimiters.first;
				} else if (i === len - 1) {
					delimiter = delimiters.last;
				} else {
					delimiter = delimiters.middle;
				}

				balloon.push(delimiter[0] + " " + pad(lines[i], maxLength) + " " + delimiter[1]);
			}

			balloon.push(" " + bottom(maxLength));
		}

		return balloon.join("\n");
	}

	function split (text, wrap) {
		text = text.replace(/\r\n?|[\n\u2028\u2029]/g, "\n").replace(/^\uFEFF/, '').replace(/\t/g, '        ');

		var lines = [];
		if (!wrap) {
			lines = text.split("\n");
		} else {
			var start = 0;
			while (start < text.length) {
				var nextNewLine = text.indexOf("\n", start);

				var wrapAt = Math.min(start + wrap, nextNewLine === -1 ? text.length : nextNewLine);

				lines.push(text.substring(start, wrapAt));
				start = wrapAt;

				// Ignore next new line
				if (text.charAt(start) === "\n") {
					start += 1;
				}
			}
		}

		return lines;
	}

	function max (lines) {
		var max = 0;
		for (var i = 0, len = lines.length; i < len; i += 1) {
			if (stringWidth(lines[i]) > max) {
				max = stringWidth(lines[i]);
			}
		}

		return max;
	}

	function pad (text, length) {
		return text + (new Array(length - stringWidth(text) + 1)).join(" ");
	}

	function top (length) {
		return new Array(length + 3).join("_");
	}

	function bottom (length) {
		return new Array(length + 3).join("-");
	}

	var balloon = {
		say: say,
		think: think
	};

	var replacer = function (cow, variables) {
		var eyes = escapeRe(variables.eyes);
		var eyeL = eyes.charAt(0);
		var eyeR = eyes.charAt(1);
		var tongue = escapeRe(variables.tongue);

		if (cow.indexOf("$the_cow") !== -1) {
			cow = extractTheCow(cow);
		}

		return cow
			.replace(/\$thoughts/g, variables.thoughts)
			.replace(/\$eyes/g, eyes)
			.replace(/\$tongue/g, tongue)
			.replace(/\$\{eyes\}/g, eyes)
			.replace(/\$eye/, eyeL)
			.replace(/\$eye/, eyeR)
			.replace(/\$\{tongue\}/g, tongue)
		;
	};

	/*
	 * "$" dollar signs must be doubled before being used in a regex replace
	 * This can occur in eyes or tongue.
	 * For example:
	 *
	 * cowsay -g Moo!
	 *
	 * cowsay -e "\$\$" Moo!
	 */
	function escapeRe (s) {
		if (s && s.replace) {
			return s.replace(/\$/g, "$$$$");
		}
		return s;
	}

	function extractTheCow (cow) {
		cow = cow.replace(/\r\n?|[\n\u2028\u2029]/g, "\n").replace(/^\uFEFF/, '');
		var match = /\$the_cow\s*=\s*<<"*EOC"*;*\n([\s\S]+)\nEOC\n/.exec(cow);

		if (!match) {
			console.error("Cannot parse cow file\n", cow);
			return cow;
		} else {
			return match[1].replace(/\\{2}/g, "\\").replace(/\\@/g, "@").replace(/\\\$/g, "$");
		}
	}

	var modes = {
		"b" : {
			eyes : "==",
			tongue : "  "
		},
		"d" : {
			eyes : "xx",
			tongue : "U "
		},
		"g" : {
			eyes : "$$",
			tongue : "  "
		},
		"p" : {
			eyes : "@@",
			tongue : "  "
		},
		"s" : {
			eyes : "**",
			tongue : "U "
		},
		"t" : {
			eyes : "--",
			tongue : "  "
		},
		"w" : {
			eyes : "OO",
			tongue : "  "
		},
		"y" : {
			eyes : "..",
			tongue : "  "
		}
	};

	var faces = function (options) {
		for (var mode in modes) {
			if (options[mode] === true) {
				return modes[mode];
			}
		}

		return {
			eyes : options.e || "oo",
			tongue : options.T || "  "
		};
	};

	var DEFAULT_COW = "$the_cow = <<\"EOC\";\n        $thoughts   ^__^\n         $thoughts  ($eyes)\\\\_______\n            (__)\\\\       )\\\\/\\\\\n             $tongue ||----w |\n                ||     ||\nEOC\n";

	var ackbar = "# Admiral Ackbar\n#\n# based on 'ack --bar' from http://beyondgrep.com/\n$the_cow = <<EOC;\n         $thoughts\n          $thoughts\n                      ?IIIIIII7II?????+\n                   ~III777II777I?+==++==+:\n                  ???I7I???I7II++=====++===\n                 ??+??????????+===~~=+++??==+\n                ??+??II??????+==~=~~=+++++==++\n               I+?????????+?+====~=~==+==++?==?\n              ?????II?????+++++=======?===~~~~==\n            ,?????II????????++++====~===::~~~~:~\n            I?I??II?+++??+?+++==~~~~:~:~:,:,,:::~\n           I??????+==+???++++=~~:~:~:,:::,:,,,,,::\n          +I?++++=+=+????+++=~~:~~:::,,,,::,,,:,:\n          I??+?+====+???+++===~~::::,::,:,,,,,,::\n         I????=~===++?+=+=~==~:~~:,,,,,,,.,,,,,:~\n        =??+?=~~~~??+?+===~~==,==~~~~,,,,..,,,.:=\n        II++==~~=++++++=~~=~,~+=?+?=I?++=..,.,,:\n     IIII?+?=====+~+++~=~~~:::=~+~===:,,,,,.,.::\n    I?=?I+??+=~=~?I?=+=~~~::,~~=~::~=::,,,,,,::\n    ?+I??++=++~,::+++~~~:::,,=~~=,~,..,::.:\n    ++=+?++~=:~::I+,~=:~,:,,,,:~~......::~,,,\n     ~=~=:.++~:,.,~=::::.,,:,.:~,:=...==~,::\n     =~?++??+=~~,.:?~.:,:,,,.,::,,~:=~=::,~\n     ++~~:~===~:~,.~::,~=~.:,..,:,,:==:.,:7\n     ~~,::...:=:,::+:~:.,~,...,.,,,,::~,,::~=\n      =~===+=~~,.::,,,:::,..,,,,,,,,,,,:,..,=+?\n      ~=~=~::~~~::,.,,,~:.+,..,,,,..,,,,...,+I?\n      ~==~:~~:~~,~=~~:,:~,:,,,,,,....,,,..+?I?I\n      ~=~=+,:~:=,:~~~~~~::::,.,,.,,.,,,..~+????I\n      ~=~==~=:~~:,~~~~~:::,::,.,,,..,,,I77I?+??II\n      +I7:::~~=~:,::~~~~.=.,~,,,,...,~7III?+??II7\n     777?+~:=~=~~:,::~~:::.,,,,,,,,,777II??I777777\n     777I==:=~::~~~~::~:::,:,:~:::,777I???777777777\n    7777+,~===~:~:~~~~:::,.~:=,,:777II???77777777777=?\n    777I~,~~~=~::~:,:,,,:=~~,,:7777I???I7777777777+=++\n  I7777I,,:,.==::::,:,,,,::::7777I+??I77777777777??I7I7,\n ,77777I::,..~~:,,,,,,.,:~I7777I+??I777777777777?I7777777,\n 77777777,...~~:,,,,,.,77777I7???II777777777777+?7777777777\n77777777777:,~~~,,=7777777I???II777777777777777+77777777777\n77777777777777777777777I+7?7II77777777777777777+777777777777\nEOC\n\n";

	var apertureBlank = "# Aperture Science logo, without the text inside\n# via http://pastebin.com/1AZwKrKp \n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n              .,-:;//;:=,\n          . :H\\@\\@\\@MM\\@M#H/.,+%;,\n       ,/X+ +M\\@\\@M\\@MM%=,-%HMMM\\@X/,\n     -+\\@MM; \\$M\\@\\@MH+-,;XMMMM\\@MMMM\\@+-\n    ;\\@M\\@\\@M- XM\\@X;. -+XXXXXHHH\\@M\\@M#\\@/.\n  ,%MM\\@\\@MH ,\\@%=            .---=-=:=,.\n  =\\@#\\@\\@\\@MX .,              -%HX\\$\\$%%%+;\n =-./\\@M\\@M\\$                  .;\\@MMMM\\@MM:\n X\\@/ -\\$MM/                    .+MM\\@\\@\\@M\\$\n,\\@M\\@H: :\\@:                    . =X#\\@\\@\\@\\@-\n,\\@\\@\\@MMX, .                    /H- ;\\@M\\@M=\n.H\\@\\@\\@\\@M\\@+,                    %MM+..%#\\$.\n /MMMM\\@MMH/.                  XM\\@MH; =;\n  /%+%\\$XHH\\@\\$=              , .H\\@\\@\\@\\@MX,\n   .=--------.           -%H.,\\@\\@\\@\\@\\@MX,\n   .%MM\\@\\@\\@HHHXX\\$\\$\\$%+- .:\\$MMX =M\\@\\@MM%.\n     =XMMM\\@MM\\@MM#H;,-+HMM\\@M+ /MMMX=\n       =%\\@M\\@M#\\@\\$-.=\\$\\@MM\\@\\@\\@M; %M%=\n         ,:+\\$+-,/H#MMMMMMM\\@= =,\n               =++%%%%+/:-.\nEOC\n";

	var aperture = "# Aperture Science logo\n# via http://pastebin.com/1AZwKrKp \n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n              .,-:;//;:=,\n          . :H\\@\\@\\@MM\\@M#H/.,+%;,\n       ,/X+ +M\\@\\@M\\@MM%=,-%HMMM\\@X/,\n     -+\\@MM; \\$M\\@\\@MH+-,;XMMMM\\@MMMM\\@+-\n    ;\\@M\\@\\@M- XM\\@X;. -+XXXXXHHH\\@M\\@M#\\@/.\n  ,%MM\\@\\@MH ,\\@%=            .---=-=:=,.\n  =\\@#\\@\\@\\@MX .,      WE      -%HX\\$\\$%%%+;\n =-./\\@M\\@M\\$         DO       .;\\@MMMM\\@MM:\n X\\@/ -\\$MM/        WHAT        .+MM\\@\\@\\@M\\$\n,\\@M\\@H: :\\@:         WE         . =X#\\@\\@\\@\\@-\n,\\@\\@\\@MMX, .        MUST        /H- ;\\@M\\@M=\n.H\\@\\@\\@\\@M\\@+,      BECAUSE       %MM+..%#\\$.\n /MMMM\\@MMH/.       WE         XM\\@MH; =;\n  /%+%\\$XHH\\@\\$=     CAN      , .H\\@\\@\\@\\@MX,\n   .=--------.           -%H.,\\@\\@\\@\\@\\@MX,\n   .%MM\\@\\@\\@HHHXX\\$\\$\\$%+- .:\\$MMX =M\\@\\@MM%.\n     =XMMM\\@MM\\@MM#H;,-+HMM\\@M+ /MMMX=\n       =%\\@M\\@M#\\@\\$-.=\\$\\@MM\\@\\@\\@M; %M%=\n         ,:+\\$+-,/H#MMMMMMM\\@= =,\n               =++%%%%+/:-.\nEOC\n";

	var armadillo = "# armadillo\n#\n# based on http://ascii.co.uk/art/armadillo\n$the_cow = <<EOC;\n         $thoughts\n          $thoughts\n               ,.-----__\n            ,:::://///,:::-.\n           /:''/////// ``:::`;/|/\n          /'   ||||||     :://'`\\\\\n        .' ,   ||||||     `/(  e \\\\\n  -===~__-'\\\\__X_`````\\\\_____/~`-._ `.\n              ~~        ~~       `~-'\nEOC\n\n";

	var atat = "# ATAT\n# from http://www.asciiworld.com/-Robots,24-.html (accessed 4/30/2014)\n$the_cow = <<EOC;\n  $thoughts                         ________\n   $thoughts                    _.-Y  |  |  Y-.,_\n    $thoughts                .-\"   |  |  |  ||   \"~-.      \n          _____     |\"\"[]\"|\" !\"\"! \"|\"==\"\" \"I      \n       .-\"{-. \"I----]_   :|------..| []  __L      \n      P-=}=(r\\_I]_[L__] _l|______l |..  |___I     \n      ^-=\\[_c=-'  ~j______[________]_L______L]    \n                    [_L--.\\_==I|I==/.--.j_I_/     \n                      j)==([\"-----`])==((_]       \n                       I--I\"~~\"\"\"~~\"I--I          \n                       |[]|         |[]|          \n                       j__l         j__l          \n                       |!!|         |!!|          \n                       |..|         |..|         \n                       )[](         )[](          \n                       ]--[         ]--[          \n                       [L_]         [L_]          \n                      /|..|\\       /|..|\\         \n                     '={--}=`     '={--}=`        \n                    .-^-r--^-.   .-^-r--^-.       \n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\nModified ATAT from Row  (the Ascii-Wizard of Oz)\nEOC\n";

	var atom = "# atom\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n                  =/;;/-\n                 +:    //\n                /;      /;\n               -X        H.\n .//;;;:;;-,   X=        :+   .-;:=;:;%;.\n M-       ,=;;;#:,      ,:#;;:=,       ,\\@\n :%           :%.=/++++/=.\\$=           %=\n  ,%;         %/:+/;,,/++:+/         ;+.\n    ,+/.    ,;\\@+,        ,%H;,    ,/+,\n       ;+;;/= \\@.  .H##X   -X :///+;\n       ;+=;;;.\\@,  .XM\\@\\$.  =X.//;=%/.\n    ,;:      :\\@%=        =\\$H:     .+%-\n  ,%=         %;-///==///-//         =%,\n ;+           :%-;;;:;;;;-X-           +:\n \\@-      .-;;;;M-        =M/;;;-.      -X\n  :;;::;;-.    %-        :+    ,-;;-;:==\n               ,X        H.\n                ;/      %=\n                 //    +;\n                  ,////,\n\nEOC\n";

	var awesomeFace = "# awesome face\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n                  \\#[/[#:xxxxxx:#[/[\\\\x\n             [/\\\\ &3N            W3& \\\\/[x\n          [[x\\@W                      W\\@x[[\\\\\n        /#&N                             N_#\n      /#\\@                                  \\@#/x\n    [/ NH_  ^\\@W               Nd_  ^\\@p      N /#\n   [[d\\@#_ zz\\@[/x3           3x:d9zz \\\\/#_N     d[[\n  /[3^[JMMMJ/////&         ^#NMMMMM ////#W     H[[\n [/\\@p/NMMMML\\@#[:^/3       d/JMMMMMMEx[# x\\\\      &/#\n /x &/LMMMMMMMMMM[_       x:MMMMMMMMMMMM /p      :/\n[/d d/ELLLLLLLLLD/&        \\#LLLLLLLLLLLL3/N      d/[\n//N   xxxxxxxxxxxxN       Wxxxxxxxxxxxxxx_       W//\n/[                                                //\n//N   p333333333333333333333333333333333p        W//\n[/d   _^/#\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\/H       \\@/[\n /:     \\\\#                              [x       :/\n [/\\@    d/x                             \\#:      &/#\n  [[H    ^[x                            [      H[[\n   [[d    _[x            &Hppp3d_      \\#\\\\N    \\@[[\n    [/ N   d#\\\\        &NzDDDDDDDDJp^ x[xN   N /#\n      /#&   N [:     pDDDDDDDDDDDDJ&#:H    &#/\n       :/#_W  W^##x 3DDDDDDDDDJN&:\\\\^p   W_#/\n          [[x&W  p& xx ^^^^ x:x \\@W   W&x/[\n             [/# &HW   WWWWN    WH& \\#/[\n                 [/[#\\\\xxxxxx\\\\#[/[\\\\x^\\@\nEOC\n";

	var banana = "# Banana \n#  http://www.ascii-art.de/ascii/ab/banana.txt\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n\n     \".           ,#  \n     \\\\ `-._____,-'=/\n  ____`._ ----- _,'_____PhS\n         `-----'\nEOC\n";

	var bearface = "##\n## acsii picture from http://www.ascii-art.de/ascii/ab/bear.txt\n##\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n     .--.              .--.\n    : (\\\\ \". _......_ .\" /) :\n     '.    `        `    .'\n      /'   _        _   `\\\\\n     /     $eye}      {$eye     \\\\\n    |       /      \\\\       |\n    |     /'        `\\\\     |\n     \\\\   | .  .==.  . |   /\n      '._ \\\\.' \\\\__/ './ _.'\n      /  ``'._-''-_.'``  \\\\\nEOC\n";

	var beavis_zen = "##\n## Beavis, with Zen philosophy removed.\n##\n$the_cow = <<EOC;\n   $thoughts         __------~~-,\n    $thoughts      ,'            ,\n          /               \\\\\n         /                :\n        |                  '\n        |                  |\n        |                  |\n         |   _--           |\n         _| =-.     .-.   ||\n         $eye|/$eye/       _.   |\n         /  ~          \\\\ |\n       (____\\@)  ___~    |\n          |_===~~~.`    |\n       _______.--~     |\n       \\\\________       |\n                \\\\      |\n              __/-___-- -__\n             /            _ \\\\\nEOC\n";

	var bees = "# Bees/beehive\n#  http://www.asciiworld.com/-Bees-.html\n$the_cow = <<EOC;\n          $thoughts\n           $thoughts\n\n\n      ^^      .-=-=-=-.  ^^\n  ^^        (`-=-=-=-=-`)         ^^\n          (`-=-=-=-=-=-=-`)  ^^         ^^\n    ^^   (`-=-=-=-=-=-=-=-`)   ^^                            ^^\n        ( `-=-=-=-(@)-=-=-` )      ^^\n        (`-=-=-=-=-=-=-=-=-`)  ^^\n        (`-=-=-=-=-=-=-=-=-`)              ^^\n        (`-=-=-=-=-=-=-=-=-`)                      ^^\n        (`-=-=-=-=-=-=-=-=-`)  ^^\n         (`-=-=-=-=-=-=-=-`)          ^^\n          (`-=-=-=-=-=-=-`)  ^^                 ^^\n      jgs   (`-=-=-=-=-`)\n             `-=-=-=-=-`\nEOC\n";

	var billTheCat = "# Bill the Cat\n#\n# Based on 'ack --th[pt]+t+'\n#  from http://beyondgrep.com/ack-2.14-single-file\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n _   /|\n \\\\'o.O'\n =(___)=\n    U\nEOC\n";

	var biohazard = "# biohazard symbol\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n              =+\\$HM####\\@H%;,\n           /H###############M\\$,\n           ,\\@################+\n            .H##############+\n              X############/\n               \\$##########/\n                %########/\n                 /X/;;+X/\n \n                  -XHHX-\n                 ,######,\n \\#############X  .M####M.  X#############\n \\##############-   -//-   -##############\n X##############%,      ,+##############X\n -##############X        X##############-\n  %############%          %############%\n   %##########;            ;##########%\n    ;#######M=              =M#######;\n     .+M###\\@,                ,\\@###M+.\n        :XH.                  .HX:\n\nEOC\n";

	var bishop = "# Bishop (Chess piece)\n#\n# from http://www.chessvariants.org/d.pieces/ascii.html\n#   by David Moeser\n#\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n    <>_\n  (\\\\)  )\n   \\\\__/\n  (____)\n   |  |\n   |__|\n  /____\\\\\n (______)\nEOC\n";

	var blackMesa = "# Black Mesa logo\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n           .-;+\\$XHHHHHHX\\$+;-.\n        ,;X\\@\\@X%/;=----=:/%X\\@\\@X/,\n      =\\$\\@\\@%=.              .=+H\\@X:\n    -XMX:                      =XMX=\n   /\\@\\@:                          =H\\@+\n  %\\@X,                            .\\$\\@\\$\n +\\@X.                               \\$\\@%\n-\\@\\@,                                .\\@\\@=\n%\\@%                                  +\\@\\$\nH\\@:                                  :\\@H\nH\\@:         :HHHHHHHHHHHHHHHHHHX,    =\\@H\n%\\@%         ;\\@M\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@H-   +\\@\\$\n=\\@\\@,        :\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@= .\\@\\@:\n +\\@X        :\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@M\\@\\@\\@\\@\\@\\@:%\\@%\n  \\$\\@\\$,      ;\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@M\\@\\@\\@\\@\\@\\@\\$.\n   +\\@\\@HHHHHHH\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@+\n    =X\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@X=\n      :\\$\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@M\\@\\@\\@\\@\\$:\n        ,;\\$\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@X/-\n           .-;+\\$XXHHHHHX\\$+;-.\nEOC\n";

	var bong = "##\n## A cow with a bong, from lars@csua.berkeley.edu\n##\n$the_cow = <<EOC;\n         $thoughts\n          $thoughts\n            ^__^ \n    _______/($eyes)\n/\\\\/(       /(__)\n   | W----|| |~|\n   ||     || |~|  ~~\n             |~|  ~\n             |_| o\n             |#|/\n            _+#+_\nEOC\n";

	var box = "# Box\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n         __________________\n        /\\\\  ______________ \\\\\n       /::\\\\ \\\\ZZZZZZZZZZZZ/\\\\ \\\\\n      /:/\\\\.\\\\ \\\\        /:/\\\\:\\\\ \\\\\n     /:/Z/\\\\:\\\\ \\\\      /:/Z/\\\\:\\\\ \\\\\n    /:/Z/__\\\\:\\\\ \\\\____/:/Z/  \\\\:\\\\ \\\\\n   /:/Z/____\\\\:\\\\ \\\\___\\\\/Z/    \\\\:\\\\ \\\\\n   \\\\:\\\\ \\\\ZZZZZ\\\\:\\\\ \\\\ZZ/\\\\ \\\\     \\\\:\\\\ \\\\\n    \\\\:\\\\ \\\\     \\\\:\\\\ \\\\ \\\\:\\\\ \\\\     \\\\:\\\\ \\\\\n     \\\\:\\\\ \\\\     \\\\:\\\\ \\\\_\\\\;\\\\_\\\\_____\\\\;\\\\ \\\\\n      \\\\:\\\\ \\\\     \\\\:\\\\_________________\\\\\n       \\\\:\\\\ \\\\    /:/ZZZZZZZZZZZZZZZZZ/\n        \\\\:\\\\ \\\\  /:/Z/    \\\\:\\\\ \\\\  /:/Z/\n         \\\\:\\\\ \\\\/:/Z/      \\\\:\\\\ \\\\/:/Z/\n          \\\\:\\\\/:/Z/________\\\\;\\\\/:/Z/\n           \\\\::/Z/_______itz__\\\\/Z/\n            \\\\/ZZZZZZZZZZZZZZZZZ/\nEOC\n";

	var brokenHeart = "# broken heart\n# via http://pastebin.com/1AZwKrKp\n# TODO: replace \"thoughts\" with \"feelings\"\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n                          .,---.\n                        ,/XM#MMMX;,\n                      -%##########M%,\n                     -\\@######%  \\$###\\@=\n      .,--,         -H#######\\$   \\$###M:\n   ,;\\$M###MMX;     .;##########\\$;HM###X=\n ,/\\@##########H=      ;################+\n-+#############M/,      %##############+\n%M###############=      /##############:\nH################      .M#############;.\n\\@###############M      ,\\@###########M:.\nX################,      -\\$=X#######\\@:\n/\\@##################%-     +######\\$-\n.;##################X     .X#####+,\n .;H################/     -X####+.\n   ,;X##############,       .MM/\n      ,:+\\$H\\@M#######M#\\$-    .\\$\\$=\n           .,-=;+\\$\\@###X:    ;/=.\n                  .,/X\\$;   .::,\n                      .,    ..\nEOC\n";

	var budFrogs = "##\n## The Budweiser frogs\n##\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n          oO)-.                       .-(Oo\n         /__  _\\\\                     /_  __\\\\\n         \\\\  \\\\(  |     ()~()         |  )/  /\n          \\\\__|\\\\ |    (-___-)        | /|__/\n          '  '--'    ==`-'==        '--'  '\nEOC\n";

	var bunny = "##\n## A cute little wabbit\n##\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts   \\\\\n        \\\\ /\\\\\n        ( )\n      .( o ).\nEOC\n";

	var C3PO = "# C3PO\n#\n# adapted from 'telnet -e x towel.blinkenlights.nl'\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n       /~\\\\\n      |oo )\n      _\\\\=/_\n     /     \\\\\n    //|/.\\\\|\\\\\\\\\n   ||  \\\\_/  ||\n   || |\\\\ /| ||\n    \\# \\\\_ _/  \\#\n      | | |\n      | | |\n      []|[]\n      | | |\n     /_]_[_\\\\\nEOC\n";

	var cake = "# Cake, from Portal \n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n            ,:/+/-\n            /M/              .,-=;//;-\n       .:/= ;MH/,    ,=/+%\\$XH@MM#@:\n      -\\$##@+\\$###@H@MMM#######H:.    -/H#\n .,H@H@ X######@ -H#####@+-     -+H###@X\n  .,@##H;      +XM##M/,     =%@###@X;-\nX%-  :M##########$.    .:%M###@%:\nM##H,   +H@@@$/-.  ,;\\$M###@%,          -\nM####M=,,---,.-%%H####M\\$:          ,+@##\n@##################@/.         :%H##@\\$-\nM###############H,         ;HM##M\\$=\n\\#################.    .=\\$M##M\\$=\n\\#################H..;XM##M\\$=          .:+\nM###################@%=           =+@MH%\n@################M/.          =+H#X%=\n=+M##############M,       -/X#X+;.\n  .;XM##########H=    ,/X#H+:,\n     .=+HM######M+/+HM@+=.\n         ,:/%XM####H/.\n              ,.:=-.\nEOC\n";

	var cakeWithCandles = "# cake with candles\n# via http://chris.com/ascii/index.php?art=events/birthday\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n       $thoughts\n                                    (\n                       (\n               )                    )             (\n                       )           (o)    )\n               (      (o)    )     ,|,            )\n              (o)     ,|,          |~\\\\    (      (o)\n              ,|,     |~\\\\    (     \\\\ |   (o)     ,|,\n              \\\\~|     \\\\ |   (o)    |`\\\\   ,|,     |~\\\\\n              |`\\\\     |`\\\\\\@\\@\\@,|,\\@\\@\\@\\@\\\\ |\\@\\@\\@\\\\~|     \\\\ |\n              \\\\ | o\\@\\@\\@\\\\ |\\@\\@\\@\\\\~|\\@\\@\\@\\@|`\\\\\\@\\@\\@|`\\\\\\@\\@\\@o |`\\\\\n             o|`\\\\\\@\\@\\@\\@\\@|`\\\\\\@\\@\\@|`\\\\\\@\\@\\@\\@\\\\ |\\@\\@\\@\\\\ |\\@\\@\\@\\@\\@\\\\ |o\n           o\\@\\@\\\\ |\\@\\@\\@\\@\\@\\\\ |\\@\\@\\@\\\\ |\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@|`\\\\\\@\\@\\@\\@\\@|`\\\\\\@\\@o\n          \\@\\@\\@\\@|`\\\\\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@|`\\\\\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\\\ |\\@\\@\\@\\@\\@\\\\ |\\@\\@\\@\\@\n          p\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\\\ |\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@|`\\\\\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@q\n          \\@\\@o\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@|`\\\\\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@o\\@\\@\n          \\@:\\@\\@\\@o\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@o\\@\\@::\\@\n          ::\\@\\@::\\@\\@o\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@o\\@\\@:\\@\\@::\\@\n          ::\\@\\@::\\@\\@\\@\\@::oo\\@\\@\\@\\@oo\\@\\@\\@\\@\\@ooo\\@\\@\\@\\@\\@o:::\\@\\@\\@::::::\n          %::::::\\@::::::\\@\\@\\@\\@:::\\@\\@\\@:::::\\@\\@\\@\\@:::::\\@\\@:::::%\n          %%::::::::::::\\@\\@::::::\\@:::::::\\@\\@::::::::::::%%\n          ::%%%::::::::::\\@::::::::::::::\\@::::::::::%%%::\n        .#::%::%%%%%%:::::::::::::::::::::::::%%%%%::%::#.\n      .###::::::%%:::%:%%%%%%%%%%%%%%%%%%%%%:%:::%%:::::###.\n    .#####::::::%:::::%%::::::%%%%:::::%%::::%::::::::::#####.\n   .######`:::::::::::%:::::::%:::::::::%::::%:::::::::\\'######.\n   .#########``::::::::::::::::::::::::::::::::::::\\'\\'#########.\n   `.#############```::::::::::::::::::::::::\\'\\'\\'#############.\\'\n    `.######################################################.\\'\n      ` .###########,._.,,,. \\#######<_\\\\##################. \\'\n         ` .#######,;:      `,/____,__`\\\\_____,_________,_____\n            `  .###;;;`.   _,;>-,------,,--------,----------\\'\n                `  `,;\\' ~~~ ,\\'\\\\######_/\\'#######  .  \\'\n                    \\'\\'~`\\'\\'\\'\\'    -  .\\'/;  -    \\'       -Catalyst\nEOC\n";

	var cat2 = "#\n#\tCat picture by Joan Stark\n#\tTransformed into cowfile by Myroslav Golub\n#\n$the_cow = <<EOC;\n       $thoughts  \n        $thoughts\n         $thoughts\n          $thoughts\n          |\\\\___/|\n         =) $eyeY$eye (=            \n          \\\\  ^  /\n           )=*=(       \n          /     \\\\\n          |     |\n         /| | | |\\\\\n         \\\\| | |_|/\\\\\n         //_// ___/\n             \\\\_) \nEOC\n";

	var cat = "# Cat\n#\n# used https://github.com/paulkaefer/flipFile.py\n#  python flipFile.py cat \" \"\n# and \n#  cat cat_flipped | sed 's/\\\\/\\\\\\\\/g' > cat.cow\n#\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts                       _\n                          / )      \n                         / /       \n      //|                \\\\ \\\\       \n   .-`^ \\\\   .-`````-.     \\\\ \\\\      \n o` {|}  \\\\_/         \\\\    / /      \n '--,  _ //   .---.   \\\\  / /       \n   ^^^` )/  ,/     \\\\   \\\\/ /        \n        (  /)      /\\\\/   /         \n        / / (     / (   /          \n    ___/ /) (  __/ __\\\\ (           \n   (((__)((__)((__(((___)          \nEOC\n\n";

	var catfence = "#\n#\tCat picture by Joan Stark\n#\tTransformed into cowfile by Myroslav Golub\n#\n$the_cow = <<EOC;\n       $thoughts     *     ,MMM8&&&.            *\n                  MMMM88&&&&&    .\n        $thoughts        MMMM88&&&&&&&\n     *           MMM88&&&&&&&&\n         $thoughts       MMM88&&&&&&&&\n                 'MMM88&&&&&&'\n          $thoughts        'MMM8&&&'      *\n          |\\\\___/|\n         =) $eyeY$eye (=            .              '\n          \\\\  ^  /\n           )=*=(       *\n          /     \\\\\n          |     |\n         /| | | |\\\\\n         \\\\| | |_|/\\\\\n  _/\\\\_/\\\\_//_// ___/\\\\_/\\\\_/\\\\_/\\\\_/\\\\_/\\\\_/\\\\_/\\\\_/\\\\_\n  |  |  |  | \\\\_) |  |  |  |  |  |  |  |  |  |\n  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |\n  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |     \n  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |\n  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |\n\nEOC\n";

	var charizardvice = "$the_cow = <<\"EOC\";\n                        $thoughts\n                         $thoughts     ___.\n                          $thoughts    L._, \\\\\n               _.,         $thoughts   <  <\\\\                _\n             ,' '           $thoughts  `.   | \\\\            ( `\n          ../, `.            $thoughts  |    .\\\\`.           \\\\ \\\\_\n         ,' ,..  .           _.,'    ||\\\\l            )  '\".\n        , ,'   \\\\           ,'.-.`-._,'  |           .  _._`.\n      ,' /      \\\\ \\\\        `' ' `--/   | \\\\          / /   ..\\\\\n    .'  /        \\\\ .         |\\\\__ - _ ,'` `        / /     `.`.\n    |  '          ..         `-...-\"  |  `-'      / /        . `.\n    | /           |L__           |    |          / /          `. `.\n   , /            .   .          |    |         / /             ` `\n  / /          ,. ,`._ `-_       |    |  _   ,-' /               ` \\\\\n / .           \\\\\"`_/. `-_ \\\\_,.  ,'    +-' `-'  _,        ..,-.    \\\\`.\n  '         .-f    ,'   `    '.       \\\\__.---'     _   .'   '     \\\\ \\\\\n' /          `.'    l     .' /          \\\\..      ,_|/   `.  ,'`     L`\n|'      _.-\"\"` `.    \\\\ _,'  `            \\\\ `.___`.'\"`-.  , |   |    | \\\\\n||    ,'      `. `.   '       _,...._        `  |    `/ '  |   '     .|\n||  ,'          `. ;.,.---' ,'       `.   `.. `-'  .-' /_ .'    ;_   ||\n|| '              V      / /           `   | `   ,'   ,' '.    !  `. ||\n||/            _,-------7 '              . |  `-'    l         /    `||\n |          ,' .-   ,' ||               | .-.        `.      .'     ||\n `'        ,'    `\".'    |               |    `.        '. -.'       `'\n          /      ,'      |               |,'    \\\\-.._,.'/'\n          .     /        .               .       \\\\    .''\n        .`.    |         `.             /         :_,'.'\n          \\\\ `...\\\\   _     ,'-.        .'         /_.-'\n           `-.__ `,  `'   .  _.>----''.  _  __  /\n                .'        /\"'          |  \"'   '_\n               /_|.-'\\\\ ,\".             '.'`__'-( \\\\\n                 / ,\"'\"\\\\,'               `/  `-.|\" m\nEOC\n";

	var charlie = "##\n## KMB is God.\n##\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n    $thoughts     ,, ＿\n        ／      ｀､\n       /   (_ﾉL_） ヽ\n      /   ´・ ・｀  l\n    （l      し     l）\n      l     ＿＿    l\n      >  ､ _      ィ\n    ／        ￣    ヽ\n   /  |              iヽ\n   |＼|              |/|\n   |  ||/＼／＼／＼/ | |\nEOC\n";

	var cheese = "##\n## The cheese from milk & cheese\n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n      _____   _________\n     /     \\\\_/         |\n    |                 ||\n    |                 ||\n   |    ###\\\\  /###   | |\n   |     $eye  \\\\/  $eye    | |\n  /|                 | |\n / |        <        |\\\\ \\\\\n| /|                 | | |\n| |     \\\\_______/   |  | |\n| |        $tongue       | / /\n/||                 /|||\n   ----------------|\n        | |    | |\n        ***    ***\n       /___\\\\  /___\\\\\nEOC\n";

	var chessmen = "# Chessmen Lineup\n#\n# based on ASCII chess pieces from http://www.chessvariants.org/d.pieces/ascii.html\n#\n# used https://github.com/paulkaefer/connectFiles.py\n#   to \"glue\" the pieces together into one file\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts \n      $thoughts\n       $thoughts\n                                           .::.                      \n                                           _::_                      \n                                 ()      _/____\\\\_                    \n                               <~~~~>    \\\\      /                    \n                       <>_      \\\\__/      \\\\____/      <>_            \n           __/\"\"\"\\\\   (\\\\)  )    (____)     (____)    (\\\\)  )   __/\"\"\"\\\\ \n  WWWWWW  ]___ 0  }   \\\\__/      |  |       |  |      \\\\__/   ]___ 0  }  WWWWWW\n   |  |       /   }  (____)     |  |       |__|     (____)      /   }   |  |\n   |  |     /~    }   |  |      |__|      /    \\\\     |  |     /~    }   |  |\n   |__|     \\\\____/    |__|     /____\\\\    (______)    |__|     \\\\____/    |__|\n  /____\\\\    /____\\\\   /____\\\\   (______)  (________)  /____\\\\    /____\\\\   /____\\\\\n (______)  (______) (______) (________) /________\\\\ (______)  (______) (______)\n\n    __        __       __        __         __        __        __       __\n   (  )      (  )     (  )      (  )       (  )      (  )      (  )     (  )\n    ||        ||       ||        ||         ||        ||        ||       ||\n   /__\\\\      /__\\\\     /__\\\\      /__\\\\       /__\\\\      /__\\\\      /__\\\\     /__\\\\\n  (____)    (____)   (____)    (____)     (____)    (____)    (____)   (____)\nEOC\n";

	var chito = "#\n# ちーちゃん\n#\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n                -一     一-\n        ／                       ＼\n       /             ________\n      /     -~                     ミ､\n      レ'     _  一ｧiァ ￢}￣Tii一- _  ＼\n    ／    --::|::::/斗士  /   |[_Vい＿＞」\n  ／ イ「::::|:::Y/  ｲ::ハ      ｨ-ﾐヽい\n  ＜___｜:::へ|::|{ 乂-夕     {::ｄﾘ|い\n        ＼八 |::｜             `''   ﾊ|\n    ＿ --＼ヽ|::|                  .ｲ ﾘ\n  ／------.ゝ|:ﾄ|        -       ィ:|\n  ＼        ＞ミ|`ヽ!ﾆ  T  ﾌ￣.≧｜:/\n     ∨         |::\\/ }-/く＼   /｜/ \nEOC\n";

	var clawArm = "# claw arm\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n       X MM X\n       X MM X\n       X MM X\n       X MM X\n       + HX +\n     ,=\\$\\$XX%/-\n   =X#########\\@%-\n  ;##############=\n -###############M,\n ;##\\@\\@\\@######M\\@###=\n .+:;+:=H##\\$=:/:;H.\n - +###- \\## :###,,;\n +\\@:/%;-H##H==/::H;\n  /#\\@/-=+\\$\\$%::+H#\\$\n  \\$#%-,      ,.:##-\n -\\@/            =X%.\n %H=             -\\$;\n  =HH,         .%M;\n   /MM/       :\\@M/.\n    .:XX,   -\\$H:.\nEOC\n";

	var clippy = "# Clippy\n#\n# from http://www.reddit.com/r/commandline/comments/2lb5ij/what_is_your_favorite_ascii_art/cltg01p\n#\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n     __ \n    /  \\\\  \n    |  |\n    @  @\n    |  |\n    || |/ \n    || || \n    |\\\\_/|\n    \\\\___/\nEOC\n\n";

	var companionCube = "# Companion Cube from Portal\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n\n +\\@##########M/             :\\@#########\\@/\n \\##############\\$;H#######\\@;+#############\n \\###############M########################\n \\##############X,-/++/+%+/,%#############\n \\############M\\$:           -X############\n \\##########H;.      ,--.     =X##########\n :X######M;     -\\$H\\@M##MH%:    :H#######\\@\n   =%#M+=,   ,+\\@#######M###H:    -=/M#%\n   %M##\\@+   .X##\\$, ./+- ./###;    +M##%\n   %####M.  /###=         \\@##M.   X###%\n   %####M.  ;M##H:.     =\\$###X.   \\$###%\n   %####\\@.   /####M\\$-./\\@#####:    %###%\n   %H#M/,     /H###########\\@:     ./M#%\n  ;\\$H##\\@\\@H:    .;\\$HM#MMMH\\$;,   ./H\\@M##M\\$=\n X#########%.      ..,,.     .;\\@#########\n \\###########H+:.           ./\\@###########\n \\##############/ ./%%%%+/.-M#############\n \\##############H\\$\\@#######\\@\\@##############\n \\##############X%########M\\$M#############\n +M##########H:            .\\$##########X=\nEOC\n";

	var cower = "##\n## A cowering cow\n##\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n        ,__, |    | \n        ($eyes)\\\\|    |___\n        (__)\\\\|    |   )\\\\_\n         $tongue  |    |_w |  \\\\\n             |    |  ||   *\n\n             Cower....\nEOC\n";

	var cowfee = "$the_cow = <<EOC;\n   $thoughts      {\n    $thoughts  }   }   {\n      {   {  }  }\n       }   }{  {\n      {  }{  }  }\n     ( }{ }{  { )\n    .-{   }   }-.\n   ( ( } { } { } )\n   |`-.._____..-'|\n   |             ;--.\n   |   (__)     (__  \\\\\n   |   ($eyes)      | )  )\n   |    \\\\/       |/  /\n   |     $tongue      /  /\n   |            (  /\n   \\\\             y'\n    `-.._____..-'\nEOC\n";

	var cthulhuMini = "# Cthulhu\n#\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n\n      ^(;,;)^\n\nEOC\n\n";

	var cube = "# Cube\n#\n# from http://www.reddit.com/r/commandline/comments/2lb5ij/what_is_your_favorite_ascii_art/cltrase\n#   also available at https://gist.github.com/th3m4ri0/6e3f631866da31d05030\n# \n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n       ____________\n      /\\\\  ________ \\\\\n     / /\\\\ \\\\______/\\\\ \\\\\n    / / /\\\\ \\\\  / /\\\\ \\\\ \\\\\n   / / /__\\\\ \\\\/ / /\\\\ \\\\ \\\\\n  / /_/____\\\\ \\\\/_/__\\\\_\\\\ \\\\\n  \\\\ \\\\ \\\\____/ / ________ \\\\\n   \\\\ \\\\ \\\\  / / /\\\\ \\\\  / / /\n    \\\\ \\\\ \\\\/ / /\\\\ \\\\ \\\\/ / /\n     \\\\ \\\\/ / /__\\\\_\\\\/ / /\n      \\\\  / /______\\\\/ /\n       \\\\/___________/\nEOC\n\n";

	var daemon = "##\n## 4.4 >> 5.4\n##\n$the_cow = <<EOC;\n   $thoughts         ,        ,\n    $thoughts       /(        )`\n     $thoughts      \\\\ \\\\___   / |\n            /- _  `-/  '\n           (/\\\\/ \\\\ \\\\   /\\\\\n           / /   | `    \\\\\n           $eye $eye   ) /    |\n           `-^--'`<     '\n          (_.)  _  )   /\n           `.___/`    /\n             `-----' /\n<----.     __ / __   \\\\\n<----|====O)))==) \\\\) /====\n<----'    `--' `.__,' \\\\\n             |        |\n              \\\\       /\n        ______( (_  / \\\\______\n      ,'  ,-----'   |        \\\\\n      `--{__________)        \\\\/\nEOC\n";

	var dalek = "# Dalek\n# from http://www.ascii-art.de/ascii/def/dr_who.txt (accessed 4/30/2014)\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n              ___\n      D>=G==='   '.\n            |======|\n            |======|\n        )--/]IIIIII]\n           |_______|\n           C O O O D\n          C O  O  O D\n         C  O  O  O  D\n         C__O__O__O__D\nsnd     [_____________]\nEOC\n";

	var dalekShooting = "# Dalek\n# from http://www.asciiworld.com/-Robots,24-.html (accessed 4/30/2014)\n$the_cow = <<EOC;\n                                    $thoughts\n                                     $thoughts\n                                                         ____                   \n                                               [(=]|[==/   @  \\\\     \n                                                      |--------|                \n     *                                     *  .       ==========                \n.  / *    .                         *   .* . * /.     ==========                \n / /  .                      *   .    *  \\\\. * /      ||||||||||||               \n =-=-=-=-=-=-----==-=--=-=--=-=-=-=---=--= -. %%%%%%[-- ||||||||||              \n  \\\\  \\\\ .                             *  (===========[  /=========]              \n.  \\\\   *  *                          .    /  * \\\\   |==============]             \n         *                        *      *         C @ @ @ @ @ @ |D             \n        *  *                          .           /              |              \n                                         .       C  @ @ @  @ @  @ |D            \n          *                          *          /                 |             \n                                               C  @  @  @  @  @  @ |D           \n                                              /                    |            \n                                             C  @   @   @   @  @  @ |D          \n                                            /                       |           \n                                           |@@@@@@@@@@@@@@@@@@@@@@@@@|          \n                                            -------------------------           \nModified from howard1\\@vax.oxford.ac.uk\nEOC\n";

	var dockerWhale = "##\n## docker whale\n##\n$the_cow = <<EOC;\n         $thoughts\n          $thoughts\n                    ##        .\n              ## ## ##       ==\n           ## ## ## ##      ===\n       /\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\\___/ ===\n  ~~~ {~~ ~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~\n       \\______ o          __/\n         \\    \\        __/\n          \\____\\______/\n\nEOC\n";

	var doge = "##\n## Doge\n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n\n           _                _\n          / /.           _-//\n         / ///         _-   /\n        //_-//=========     /\n      _///        //_ ||   ./\n    _|                 -__-||\n   |  __              - \\\\   \\\n  |  |#-       _-|_           |\n  |            |#|||       _   |  \n |  _==_                       ||\n- ==|.=.=|_ =                  |\n|  |-|-  ___                  |\n|    --__   _                /\n||     ===                  |\n |                     _. //\n  ||_         __-   _-  _|\n     \\_______/  ___/  _|\n                   --*\nEOC\n";

	var dolphin = "# dolphin (tiny)\n#\n# from http://www.chris.com/ascii/index.php?art=animals/other%20(water)\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n               ,\n             __)\\\\_  \n       (\\_.-'    a`-.\n  jgs  (/~~````(/~^^` \n\nEOC\n";

	var dragonAndCow = "##\n## A dragon smiting a cow, possible credit to kube@csua.berkeley.edu\n##\n$the_cow = <<EOC;\n                       $thoughts                    ^    /^\n                        $thoughts                  / \\\\  // \\\\\n                         $thoughts   |\\\\___/|      /   \\\\//  .\\\\\n                          $thoughts  /O  O  \\\\__  /    //  | \\\\ \\\\           *----*\n                            /     /  \\\\/_/    //   |  \\\\  \\\\          \\\\   |\n                            \\@___\\@`    \\\\/_   //    |   \\\\   \\\\         \\\\/\\\\ \\\\\n                           0/0/|       \\\\/_ //     |    \\\\    \\\\         \\\\  \\\\\n                       0/0/0/0/|        \\\\///      |     \\\\     \\\\       |  |\n                    0/0/0/0/0/_|_ /   (  //       |      \\\\     _\\\\     |  /\n                 0/0/0/0/0/0/`/,_ _ _/  ) ; -.    |    _ _\\\\.-~       /   /\n                             ,-}        _      *-.|.-~-.           .~    ~\n            \\\\     \\\\__/        `/\\\\      /                 ~-. _ .-~      /\n             \\\\____($eyes)           *.   }            {                   /\n             (    (--)          .----~-.\\\\        \\\\-`                 .~\n             //__\\\\\\\\$tongue\\\\__ Ack!   ///.----..<        \\\\             _ -~\n            //    \\\\\\\\               ///-._ _ _ _ _ _ _{^ - - - - ~\nEOC\n";

	var dragon = "##\n## The Whitespace Dragon\n##\n$the_cow = <<EOC;\n      $thoughts                    / \\\\  //\\\\\n       $thoughts    |\\\\___/|      /   \\\\//  \\\\\\\\\n            /$eye  $eye  \\\\__  /    //  | \\\\ \\\\    \n           /     /  \\\\/_/    //   |  \\\\  \\\\  \n           \\@_^_\\@'/   \\\\/_   //    |   \\\\   \\\\ \n           //_^_/     \\\\/_ //     |    \\\\    \\\\\n        ( //) |        \\\\///      |     \\\\     \\\\\n      ( / /) _|_ /   )  //       |      \\\\     _\\\\\n    ( // /) '/,_ _ _/  ( ; -.    |    _ _\\\\.-~        .-~~~^-.\n  (( / / )) ,-{        _      `-.|.-~-.           .~         `.\n (( // / ))  '/\\\\      /                 ~-. _ .-~      .-~^-.  \\\\\n (( /// ))      `.   {            }                   /      \\\\  \\\\\n  (( / ))     .----~-.\\\\        \\\\-'                 .~         \\\\  `. \\\\^-.\n             ///.----..>        \\\\             _ -~             `.  ^-`  ^-_\n               ///-._ _ _ _ _ _ _}^ - - - - ~                     ~-- ,.-~\n                                                                  /.-~\nEOC\n";

	var ebi_furai = "#\n# えびフライ\n#\n\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n      ,.,,､,..,､､.,､,､､.,_          ／i\n    ;'`;、､:、..:、:,:,.::｀'::ﾞ\":,'´ --i\n    '､;:..: ,:.､.:',.:.::_.;..;:.‐'ﾞ\n\nEOC\n";

	var elephant2 = "# Elephant\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts                                 \n      /  \\\\~~~/  \\\\         \n     (    ..     )----,      \n      \\\\__     __/      \\\\     \n        )|  /)         |\\\\    \n         | /\\\\  /___\\\\   / ^   \n          \"-|__|   |__|      \nEOC\n";

	var elephant = "##\n## An elephant out and about\n##\n$the_cow = <<EOC;\n $thoughts     /\\\\  ___  /\\\\\n  $thoughts   // \\\\/   \\\\/ \\\\\\\\\n     ((    $eye $eye    ))\n      \\\\\\\\ /     \\\\ //\n       \\\\/  | |  \\\\/ \n        |  | |  |  \n        |  | |  |  \n        |   o   |  \n        | |   | |  \n        |m|   |m|  \nEOC\n";

	var elephantInSnake = "##\n## Do we need to explain this?\n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts              ....       \n           ........    .      \n          .            .      \n         .             .      \n.........              .......\n..............................\n\nElephant inside ASCII snake\nEOC\n";

	var explosion = "# Explosion\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n            .+\n             /M;\n              H#@:              ;,\n              -###H-          -@/\n               %####\\$.  -;  .%#X\n                M#####+;#H :M#M.\n..          .+/;%#########X###-\n -/%H%+;-,    +##############/\n    .:\\$M###MH\\$%+############X  ,--=;-\n        -/H#####################H+=.\n           .+#################X.\n         =%M####################H;.\n            /@###############+;;/%%;,\n         -%###################\\$.\n       ;H######################M=\n    ,%#####MH\\$%;+#####M###-/@####%\n  :\\$H%+;=-      -####X.,H#   -+M##@-\n .              ,###;    ;      =\\$##+\n                .#H,               :XH,\n                 +                   .;-\nEOC\n";

	var eyes = "##\n## Evil-looking eyes\n##\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n                                   .::!!!!!!!:.\n  .!!!!!:.                        .:!!!!!!!!!!!!\n  ~~~~!!!!!!.                 .:!!!!!!!!!UWWW\\$\\$\\$ \n      :\\$\\$NWX!!:           .:!!!!!!XUWW\\$\\$\\$\\$\\$\\$\\$\\$\\$P \n      \\$\\$\\$\\$\\$##WX!:      .<!!!!UW\\$\\$\\$\\$\"  \\$\\$\\$\\$\\$\\$\\$\\$# \n      \\$\\$\\$\\$\\$  \\$\\$\\$UX   :!!UW\\$\\$\\$\\$\\$\\$\\$\\$\\$   4\\$\\$\\$\\$\\$* \n      ^\\$\\$\\$B  \\$\\$\\$\\$\\\\     \\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$   d\\$\\$R\" \n        \"*\\$bd\\$\\$\\$\\$      '*\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$o+#\" \n             \"\"\"\"          \"\"\"\"\"\"\" \nEOC\n";

	var fatBanana = "# fatter banana\n# via https://www.reddit.com/r/cowsay/comments/3bkpwv/any_love_for_bananasay/\n$the_cow = <<EOC;\n           $thoughts\n            $thoughts\n        \"-.. __      __.='>\n         `.     \"\"\"\"\"   ,'\n           \"-..__   _.-\"\n   ~ ~~ ~ ~  ~   \"\"\"  ~~  ~\nEOC\n";

	var fatCow = "# fatter cow\n# via https://www.reddit.com/r/cowsay/comments/39htd0/with_all_this_reddit_hype_what_about_a_little/\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n\n    A__A\n   ( OO )\\\\_----__\n   (____)\\\\      )\\\\/\\\\\n        ||      |\n        ||`---w||\nEOC\n";

	var fence = "$the_cow = <<EOC;\n                          $thoughts\n                           $thoughts         __.----.___\n           ||            ||  (\\\\(__)/)-'||      ;--` ||\n          _||____________||___`($eyes)'___||______;____||_\n          -||------------||----)  (----||-----------||-\n          _||____________||___(o  o)___||______;____||_\n          -||------------||----`--'----||-----------||-\n           ||            ||     $tongue `|| ||| || ||     ||jgs\n        ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\nEOC\n";

	var fire = "# Fire\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n      $thoughts\n       $thoughts\n                     -\\$-\n                    .H##H,\n                   +######+\n                .+#########H.\n              -\\$############\\@.\n            =H###############\\@  -X:\n          .\\$##################:  \\@#\\@-\n     ,;  .M###################;  H###;\n   ;\\@#:  \\@###################\\@  ,#####:\n -M###.  M#################\\@.  ;######H\n M####-  +###############\\$   =\\@#######X\n H####\\$   -M###########+   :#########M,\n  /####X-   =########%   :M########\\@/.\n    ,;%H\\@X;   .\\$###X   :##MM\\@%+;:-\n                 ..\n  -/;:-,.              ,,-==+M########H\n -##################\\@HX%%+%%\\$%%%+:,,\n    .-/H%%%+%%\\$H\\@###############M\\@+=:/+:\n/XHX%:#####MH%=    ,---:;;;;/%%XHM,:###\\$\n\\$\\@#MX %+;-                           .\nEOC\n";

	var flamingSheep = "##\n## The flaming sheep, contributed by Geordan Rosario (geordan@csua.berkeley.edu)\n##\n$the_cow = <<EOC;\n  $thoughts            .    .     .   \n   $thoughts      .  . .     `  ,     \n    $thoughts    .; .  : .' :  :  : . \n     $thoughts   i..`: i` i.i.,i  i . \n      $thoughts   `,--.|i |i|ii|ii|i: \n           U${eyes}U\\\\.'\\@\\@\\@\\@\\@\\@`.||' \n           \\\\__/(\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@)'  \n             $tongue (\\@\\@\\@\\@\\@\\@\\@\\@)    \n                `YY~~~~YY'    \n                 ||    ||     \nEOC\n";

	var fox = "# Fox\n# http://www.retrojunkie.com/asciiart/animals/foxes.htm\n$the_cow = <<EOC;\n$thoughts\n $thoughts\n   /\\\\   /\\\\   Todd Vargo\n  //\\\\\\\\_//\\\\\\\\     ____\n  \\\\_     _/    /   /\n   / * * \\\\    /^^^]\n   \\\\_\\\\O/_/    [   ]\n    /   \\\\_    [   /\n    \\\\     \\\\_  /  /\n     [ [ /  \\\\/ _/\n    _[ [ \\\\  /_/\nEOC\n";

	var ghostbusters = "##\n## Ghostbusters!\n##\n$the_cow = <<EOC;\n          $thoughts\n           $thoughts\n            $thoughts          __---__\n                    _-       /--______\n               __--( /     \\\\ )XXXXXXXXXXX\\\\v.\n             .-XXX(   $eye   $eye  )XXXXXXXXXXXXXXX-\n            /XXX(       U     )        XXXXXXX\\\\\n          /XXXXX(              )--_  XXXXXXXXXXX\\\\\n         /XXXXX/ (      O     )   XXXXXX   \\\\XXXXX\\\\\n         XXXXX/   /            XXXXXX   \\\\__ \\\\XXXXX\n         XXXXXX__/          XXXXXX         \\\\__---->\n ---___  XXX__/          XXXXXX      \\\\__         /\n   \\\\-  --__/   ___/\\\\  XXXXXX            /  ___--/=\n    \\\\-\\\\    ___/    XXXXXX              '--- XXXXXX\n       \\\\-\\\\/XXX\\\\ XXXXXX                      /XXXXX\n         \\\\XXXXXXXXX   \\\\                    /XXXXX/\n          \\\\XXXXXX      >                 _/XXXXX/\n            \\\\XXXXX--__/              __-- XXXX/\n             -XXXXXXXX---------------  XXXXXX-\n                \\\\XXXXXXXXXXXXXXXXXXXXXXXXXX/\n                  \"\"VXXXXXXXXXXXXXXXXXXV\"\"\nEOC\n";

	var ghost = "# art by Joan G. Stark, https://en.wikipedia.org/wiki/Joan_Stark\n$the_cow = <<\"EOC\";\n     $thoughts     .-.\n      $thoughts  .'   `.\n       $thoughts :g g   :\n        $thoughts: o    `.\n        :         ``.\n       :             `.\n      :  :         .   `.\n      :   :          ` . `.\n       `.. :            `. ``;\n          `:;             `:'\n             :              `.\n              `.              `.     .\n                `'`'`'`---..,___`;.-'\nEOC\n\n";

	var glados = "# GLaDOS from Portal\n# via http://pastebin.com/1AZwKrKp \n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n       \\#+ \\@      \\# \\#              M#\\@\n .    .X  X.%##\\@;# \\#   +\\@#######X. \\@#%\n   ,==.   ,######M+  -#####%M####M-    \\#\n  :H##M%:=##+ .M##M,;#####/+#######% ,M#\n .M########=  =\\@#\\@.=#####M=M#######=  X#\n :\\@\\@MMM##M.  -##M.,#######M#######. =  M\n             \\@##..###:.    .H####. \\@\\@ X,\n   \\############: \\###,/####;  /##= \\@#. M\n           ,M## ;##,\\@#M;/M#M  \\@# X#% X#\n.%=   \\######M## \\##.M#:   ./#M ,M \\#M ,#\\$\n\\##/         \\$## \\#+;#: \\#### ;#/ M M- \\@# :\n\\#+ \\#M\\@MM###M-;M \\#:\\$#-##\\$H# .#X \\@ + \\$#. \\#\n      \\######/.: \\#%=# M#:MM./#.-#  \\@#: H#\n+,.=   \\@###: /\\@ %#,\\@  \\##\\@X \\#,-#\\@.##% .\\@#\n\\#####+;/##/ \\@##  \\@#,+       /#M    . X,\n   ;###M#\\@ M###H .#M-     ,##M  ;\\@\\@; \\###\n   .M#M##H ;####X ,\\@#######M/ -M###\\$  -H\n    .M###%  X####H  .\\@\\@MM\\@;  ;\\@#M\\@\n      H#M    /\\@####/      ,++.  / ==-,\n               ,=/:, .+X\\@MMH\\@#H  \\#####\\$=\nEOC\n";

	var goat2 = "#\n#\tCodeGoat.io: https://github.com/danyshaanan/goatsay\n#\n$the_cow = <<EOC;\n        $thoughts\n         $thoughts\n          )__(\n         '|$eyes|'________/\n          |__|         |\n           $tongue||\"\"\"\"\"\"\"||\n             ||       ||\n\nEOC\n";

	var goat = "##\n## ejm97 http://www.ascii-art.de/ascii/ghi/goat.txt\n##\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n         $thoughts  _))\n           > $eye\\\\     _~\n           `;'\\\\\\\\__-' \\\\_\n              | )  _ \\\\ \\\\\n             / / ``   w w\n            w w\nEOC\n\n\n\n\n\n";

	var goldenEagle = "# Golden Eagle (Marquette University mascot)\n# \n$the_cow = <<EOC;\n    $thoughts                                       ,:=+++++???++?+=+=\n     $thoughts                               :+?????\\$MUUUUUUUUUMO+??~\n      $thoughts                         :+I??\\$UUUUMUMMMMUUMUMUUMUUMM???I+:\n       $thoughts                     ,+??+ZOUUMMUMMMUUUUUMUUUMUUUMUMUUMZI+?+:\n        $thoughts                 ~I?+MMUMUUUMUUUOOUMMMMMMUUUUUMMMUUUUUUMMUM$??~\n                       I?+7MMMMMUUO7?+?IUMMMMMMMMUUMUUMUUUUUUUUUUMMMUMO?I\n                    ~I?+MMMUUUO????+?IOUZ7?,.......,+\\$\\$OUMUUUUUUMUMUUUMUU+I:\n                   =??\\$UMUUU7++??????II???????=.....,?OUUMMMUUMUUUUUUUMMUUU+=\n                 +??UUMMM7??????+??+?????+??=,...\\$MUMUUUUMUUMUUUUUUUUMMUM7II??=\n               ,+?IUMMMI???III?++??+?????+~....... ......MUUMUUUUUUUUUMMU7?~\n               IIUMMM+?+?IUUUUMUUM7I?????????????I?+=:......MUUMMUMUMMMMUUU+~\n             :?+UMMU+?+?7UMMUUUZ7\\$\\$7????+++????????????=.....+UMUUUUMMMMMMUZ?\n             ?+UMUM???+MMMMMU?++???????????++????????++????....OMMMUMMMMUMMUI:\n            +\\$MMUM?+?ZMMU:\\$MM???OUUU+??+???+????????????????,...UMUMUUUUMMUMM?~\n            IUUUU?I?OUUU,..UU?IMMUUMUUI???+?????????????????I,..:UUMUUUUMMUMU?+\n            ?UUUMUM\\$UMUU~..UUUUU\\$,IUUUMM7+?????????????????+?I~..UUUUUUUMMMUU+?\n            ?OUMUUUUMMUI+.?UUUU=...~UMMUU\\$?????+???????????????..MUMUUUUUMUMU??\n           :??IUUMMUMUMMOMUU7........OUUUMMU?I????????????????I..MUMUUUUMUUMU?+\n         +IIUMUO.IUUUUUUO..............?UMUMUM7??????????????+?..UUMUUUUMUUMU?=\n       ,IZMMU,.:UU7:..........,UUUMUZ....MUUMMO+???????????????..UUUUUUUUUUU?:\n       IZUUU:..UUUI=....... IUUUUUMMUZ,.MMUMU$?+???????????????.MUUMUUUUMUMUI\n     ,+IUUM..O=..........\\$UUMMMUU?~....UMMUUI?????????????????=.UMMUUUUUMMU?+\n     +?UMU~............OUMMUU~..... .UUMUMM+?????????????????=.UMMUMUMUUMOI=\n     ?\\$MU~...    ...:MUMU=~........,UUMMMUI+????????????????+IUMMMUUUUMUU?+\n     +OMU....   ...?UMU=..:~~,.....MMMUUU+?+????????????????~MMUUUUMUMMU?+~   \n     ?OMU~ .. ...?UMUUUMUMUMUMUMUUUMUUMUI???????????????+?+OUUMUUMMMMUIUUUMO,\n     ??UMU~.....\\$MUUUOM???UMMUUUMMMUUMM7?++????????????++OMMMUMUUMMUI??UMIU+~\n     :?7UUU\\$...UMMM?I~,  +?MMUUMMMMMMUU?????????????+??\\$UMMMMUUUMU\\$?: ,??I?:\n       ?IMUMUUZMU+?,      =?UMMMMMMMMO??????????????+UMUMMUMUMUU?I~\n        ?+\\$MUMUMU??        ?MMMMMMMMU??+???????????IUMMMUUUMUUZ?=\n          ,+???ZUO?~       +ZUMMUMUMU???++??+???IUMMUMUMMUUO??~\n               ,,:~=       ,?UMUMUMU???+??+?+?7UMUMUMUMUI??:\n                            ?UMUMMMM?+??++?ZUUMUMUMUZ++?,                  \n                            ?UMMMMMO+???MMUMUMUMUMOII=,                       \n                            ?UUUMUUZOMUUMUMMUMM+??=\n                            ?UMMUMMUUUMUMM\\$???~\n                           ,?UMUUMUUU\\$?+?~:                                  \n                           :IUUUM?+?I=:                                  \n                           ????~,\nEOC\n";

	var hand = "##\n## これが私の本当の姿だ！\n##  \n##\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n                           __ \n                  l^ヽ    /  }    _\n                  |  |   /  /   ／  )\n                  |  |  /  /  ／  ／ _\n                  j. し'  / ／  ／ ／  )\n                 /  .＿__ ´  ／ ／  ／\n                /   {  /:｀ヽ ｀¨ ／\n               /     ∨::::::ﾊ   ／\n              |廴     ＼:::ノ}  /\n    {￣￣￣￣ヽ  廴     ｀ー'  ー-､\n    ヽ ＿＿_   ＼ 廴        ＿＿＿ﾉ\n        ／       ＼ 辷_´￣\n      ／           ﾍ￣\n    ／             ,ﾍ\n                  /、ﾍ\n                 /＼__ﾉ\nEOC\n\n";

	var happyWhale = "# happy whale\n#\n# modified from http://www.chris.com/ascii/index.php?art=animals/other%20(water)\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n     $thoughts\n        __ \\ / __\n       /  \\\\ | /  \\\\\n           \\\\|/\n       _.---v---.,_\n      /            \\\\  /\\\\__/\\\\\n     /              \\\\ \\\\_  _/\n     |__ @           |_/ /\n      _/                / \n      \\\\       \\\\__,     /  \n   ~~~~\\\\~~~~~~~~~~~~~~`~~~\n\nEOC\n";

	var hedgehog = "##\n## A cute little hedgehog\n##\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts ..:::::::::.\n    ::::::::::::::\n   /. `::::::::::::\n  O__,_:::::::::::'\nEOC\n";

	var hellokitty = "##\n## Hello Kitty\n##\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n      /\\\\_)o<\n     |      \\\\\n     | $eye . $eye|\n      \\\\_____/\n         $tongue\nEOC\n";

	var hippie = "$the_cow = <<EOC;\n                       $thoughts              ___\n                        $thoughts            ///\\\\\\\\\\\\/----\n                         $thoughts           ||//\\\\///\\\\\\\\\\\\\\\\\n                          $thoughts         /`-.__\\\\\\\\\\\\\\\\///|\n                           $thoughts       /_  _   `--._|\n                               ___-`---.___     |\n                          ----------       `-.__|\n                       ----------( \\\\.-.$eye $eye;_  \\\\\\\\\\\\\\\\\\\\\\\\\n                      ------------| `-'-.(_)--/\\\\\\\\\\\\\\\\\\\\\n                     /////------//|   `-'       )\\\\\\\\\\\\\\\\\\\\\\\\\n                     /////------///\\\\  `--'\\\\  /\"\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n                     ////--------///\\\\  `-' /\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\   .-.  _\n                      //////------////>---'\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\  | | / )\n        _              ////////////// |__| )\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\  | |/ /\n       / `.       _    ////////.-'  >\\\\    <-._.--.\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\  _|__ /_\n      (    \\\\  . .' )    /////// ( .- (    )() ( )_)\\\\\\\\\\\\\\\\\\\\\\\\\\\\  / __)-' )\n       `-   | |/          //// ( ) ( )|--'() ( ) \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\   \\\\  `(.-')\n          .---/ _()        /// ( ) () |  /() ( )  \\\\\\\\\\\\\\\\\\\\\\\\     > ._>-'\n        ()+8 8 |            |  ( )( ) | /( ) ( )   \\\\        / \\\\/\n        ()+8 8/-()__        /  ( )( ) \\\\/ ( ) ( )\\\\   \\\\      /\\\\ /\n          |8 8|     `.     |   () ( ).--.( ) ( )-\\\\   \\\\    /   |\n        ()+||||-() (_/    _/   /| ()/ || \\\\ )  ()()\\\\   \\\\__/   /\n        .-`||||          /\\\\\\\\  / / ()|/  \\\\ ()     \\\\ `.  /|   |\n       (_  ||||        .'   _/-/  ()\\\\/||\\\\/()     \\\\-. \\\\     /\n           ||( \\\\_    .'    ( )/  ( ) `--' ( )     > ) `.  /\n        .--|_|\\\\_ \\\\ .'    .'( )_  ( )-.___.-( )  (  )    \"\"\n        `.__)-.( /.'\\\\  .'   (  )'_)-.______( ).-')'\n          (___)|  \\\\ .-'      `--'`-._.---._.(_))-'\n          (__)|| +-)'           |   /_.--.\\\\    |\n          (__)||-'              `._|`-'  ) )  _|\n            |||||                |  `.`-'.'--' /\n            |||||               .'    | |   .\\\\|\n            |||||             .'   _.-|_|     \\\\\n            |||||            /   .'.-'  \\\\\\\\     |\n            ||||||         .'     /      \\\\     \\\\\n             |||||        /     .'        \\\\     \\\\\n             |||||      .'     /           |    |\n            _|||||----./     .'            \\\\     \\\\\n         .-' |||||   `/     /               \\\\    |\n       .'     |||||   (    /                |    |\n      /       |||||   |    |\\\\                \\\\   |\n      |     .'|||||.  |    ||                |    )\n       \\\\    | |||||\\\\  |    |/                |    \\\\\n        \\\\   | ||||||  |    |                 /    |\n        |    `.||||' /     |                |     \\\\\n        |      ||||  |     \\\\                |      |\n        /      ||||| |      |\\\\             /       |\n       /       |||||_/      | \\\\            |        \\\\\n      /      ------'|       |  |           |        |\n     |      |___.---|        \\\\ |           /        |\n     |             /         | |          |         \\\\\n     |             |         \\\\/           |          |\n     |             /          |           |          |\n      \\\\           |           |           |          |\n       `.        /             \\\\          |           \\\\\n         `--.___`-_            |_         |           |\n           .-.__.-''-,_         -         |           \\\\_'\n          <`.         '.-//|-/``        (_)          _.-'\n           `._-.____.-'.|   /            '//, ,\\\\.-'`` |--.\n              `-.____.' |__/               '''\\\\      -'/ |\n                                               `.   _.// |\n                                                 `-.__.-'\n\nVK\nEOC\n";

	var hiya = "$the_cow = <<EOC;\n           $thoughts     (      )\n            $thoughts    ~(^^^^)~\n             $thoughts    ) $eyes \\\\~_          |\\\\\n              $thoughts  /     | \\\\        \\\\~ /\n                ( 0  0  ) \\\\        | |\n                 ---___/~  \\\\       | |\n                  /'__/ |   ~-_____/ |\n   o          _   ~----~      ___---~\n     O       //     |         |\n            ((~\\\\  _|         -|\n      o  O //-_ \\\\/ |        ~  |\n           ^   \\\\_ /         ~  |\n                  |          ~ |\n                  |     /     ~ |\n                  |     (       |\n                   \\\\     \\\\      /\\\\\n                  / -_____-\\\\   \\\\ ~~-*\n                  |  /       \\\\  \\\\       .==.\n                  / /         / /       |  |\n                /~  |      //~  |       |__|         W<\n                ~~~~        ~~~~\nEOC\n";

	var hiyoko = "##\n## ひよ子\n##\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n    $thoughts\n      ,､ ,._\n      ﾉ ・  ヽ\n     / :::   i  \n    / :::    ﾞ､\n   ,i:::       `ｰ-､\n   |:::           i\n   !::::..        ﾉ\n    `ー――――'\" \nEOC\n";

	var homer = "# Homer Simpson\n#\n# from http://www.reddit.com/r/textfiles/comments/2s9ybk/random_ascii_art/\n#\n$the_cow = <<EOC;\n            $thoughts\n             $thoughts                __ \n                   _ ,___,-'\",-=-. \n       __,-- _ _,-'_)_  (\"\"`'-._\\\\ `. \n    _,'  __ |,' ,-' __)  ,-     /. | \n  ,'_,--'   |     -'  _)/         `\\\\ \n,','      ,'       ,-'_,`           : \n,'     ,-'       ,(,-(              : \n     ,'       ,-' ,    _            ; \n    /        ,-._/`---'            / \n   /        (____)(----. )       ,' \n  /         (      `.__,     /\\\\ /, \n :           ;-.___         /__\\\\/| \n |         ,'      `--.      -,\\\\ | \n :        /            \\\\    .__/ \n  \\\\      (__            \\\\    |_ \n   \\\\       ,`-, *       /   _|,\\\\ \n    \\\\    ,'   `-.     ,'_,-'    \\\\ \n   (_\\\\,-'    ,'\\\\\")--,'-'       __\\\\ \n    \\\\       /  // ,'|      ,--'  `-. \n     `-.    `-/ \\\\'  |   _,'         `. \n        `-._ /      `--'/             \\\\ \n-hrr-      ,'           |              \\\\ \n          /             |               \\\\ \n       ,-'              |               / \n      /                 |             -'\nEOC\n";

	var hypno = "$the_cow =<<\"EOC\"\n  $thoughts\n     ___        _--_\n    /    -    /     \\\\\n   ( $eyes   \\\\  (    $eyes )\n   |  $eyes _;\\\\-/|  $eyes _|\n    \\\\___/######\\\\___/\\\\\n      /##############\\\\\n     /  ######   ##  #|\n    /  ##@##@##       |\n   /    ######     ##  \\\\\n <______-------___\\\\  . //_\n    |       ____  | | //# \\\\__~__\n     \\\\      $tongue    \\\\  //###  \\\\   \\\\\n      |             /\\'  ##  ##  ##\\\\   __--~--_\n       \\\\_________- /\\\\ )    ^     ##|--########\\\\\n  /--~-_\\\\________/_  |          #@##|#######Y##|\n | \\\\ `  /|       /O/ ( ###  \\')    ##/######/###/\n \\\\  \\\\  | |       --  |  ###        /LLLLL--###/\n  \\\\_ \\\\/  |            \\\\_   \\\\    ) /####_____--\n ___ /    \\\\           /     |   _-####\\\\\n(___/     -\\\\_________/     / -- |#####@@@@@@\\'_\n (__\\\\_      __,) (.___     ,/    /#####      `@@\n      | -\\\\\\\\-          //-//      @@  @@@@@.\n      | | \\\\\\\\_       _// //      @\\'       \\'@@.\n      (.)   \\\\_)    / / //                   @@@\n                  (_) (_\\'\nEOC\n";

	var ibm = "#\n# International Business Machines\n#\n\n$the_cow = << EOC;\n  $thoughts\n   $thoughts\n\n■■■■■   ■■■■■■■■     ■■■■■       ■■■■■\n■■■■■   ■■■■■■■■■■   ■■■■■■     ■■■■■■\n ■■■     ■■■   ■■■    ■■■■■■   ■■■■■■\n ■■■     ■■■■■■■■     ■■■■■■■ ■■■■■■■\n ■■■     ■■■■■■■■     ■■■ ■■■■■■■ ■■■\n ■■■     ■■■   ■■■    ■■■  ■■■■■  ■■■\n■■■■■   ■■■■■■■■■■   ■■■■   ■■■   ■■■■\n■■■■■   ■■■■■■■■     ■■■■    ■    ■■■■\nEOC\n";

	var iwashi = "##\n## いわし\n##  \n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n         ＿＿＿＿ ＿＿＿＿＿__\n      ｨ''  ＠ :. ,! ，， ， ，￣￣ ¨` ‐-            ＿＿\n       ＼    ノ   i            ’ ’’ ’’､_;:`:‐.-_-‐ニ＝=彳\n         ｀ ＜. _  .ｰ ､                       !三  ＜\n                 ｀¨  ‐= . ＿＿＿_.. ﾆ=-‐‐`'´｀ﾐ､   三＞\n                                                 ￣￣\nEOC\n";

	var jellyfish = "# jellyfish\n#\n# from http://ascii.co.uk/art/jellyfish\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n\n         .-;\\':\\':\\'-.\n        {\\'.\\'.\\'.\\'.\\'.}\n         )        \\'`.\n        \\'-. ._ ,_.-=\\'\n          `). ( `);(\n          (\\'. .)(,\\'.)\n           ) ( ,\\').(\n          ( .\\').\\'(\\').\n          .) (\\' ).(\\'\n           '  ) (  ).\n            .\\'( .)\\'\n              .).\\'\njgs\n\nEOC\n";

	var karl_marx = "##\n## Karl Marx\n##  \n##\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n                   ,―ヾヽヽ/ｖへ／⌒ー\n                , ⌒ヽ ヽ ヽ / ／ ノ  ⌒ヽ、\n              / ／ヾ,ゞ -ゞゞゞ､_ ⌒  ノ ヽ\n            ／  ／            `ヾ  ー   ミヽ\n          ,/   /                   ヾ ＼  ヽﾐ\n         /    /                      ゞ      ヽ\n         i   /                       /      ＼\n        /    -=ﾆヽ､,_  ,,,,;r;==-     ヾ  ヾミ ヽ\n        | ;: `ゞﾂヽ〉^`ヾだ'=-､_        i    彡 ヽ\n        i ,   /::::/     `'''\"\"\"        ﾉ  ゞ ヾ ヽ\n        } ;  |    人､,;-,'^            /    くヾ  ）\n        /    彡ノノノﾉﾉﾉ(((((        ／ﾍミ        /\n       /     /ﾉﾉﾉﾉﾉ,.-―ミヽヾヾヾヾヾヾ     _ノ`ｰ'\"\n      ,i          -ー‐ `ゞ           ヽ   ヽ\n      彡彡                        ミ       ヽ\n''\"\"￣彡      /   /   /   /            ミ   ﾂ＼\n      ＜    /   /   /   /        ヾ   ヾ  ノﾉﾉ\n        '―彡                         ｒー'\"\n            ヾノ人,,.r--､ノノノノノり'\"\nEOC\n";

	var kilroy = "# Kilroy\n# from http://www.ascii-art.de/ascii/jkl/kilroy.txt (accessed 8/14/2014)\n$the_cow = <<EOC;\n     $thoughts \n      $thoughts\n           ,,,\n          (0 0)\n   +---ooO-(_)-Ooo---+\n   |                 |\nEOC\n\n";

	var king = "# King (Chess piece)\n#\n# from http://www.chessvariants.org/d.pieces/ascii.html\n#   by David Moeser\n#\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n    .::.\n    _::_\n  _/____\\\\_\n  \\\\      /\n   \\\\____/\n   (____)\n    |  |\n    |__|\n   /    \\\\\n  (______)\n (________)\n /________\\\\\nEOC\n";

	var kiss = "##\n## A lovers' empbrace\n##\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n             ,;;;;;;;,\n            ;;;;;;;;;;;,\n           ;;;;;'_____;'\n           ;;;(/))))|((\\\\\n           _;;((((((|))))\n          / |_\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n     .--~(  \\\\ ~))))))))))))\n    /     \\\\  `\\\\-(((((((((((\\\\\\\\\n    |    | `\\\\   ) |\\\\       /|)\n     |    |  `. _/  \\\\_____/ |\n      |    , `\\\\~            /\n       |    \\\\  \\\\           /\n      | `.   `\\\\|          /\n      |   ~-   `\\\\        /\n       \\\\____~._/~ -_,   (\\\\\n        |-----|\\\\   \\\\    ';;\n       |      | :;;;'     \\\\\n      |  /    |            |\n      |       |            |\nEOC\n";

	var kitten = "# Kitten\n#\n# based on rfksay by Andrew Northern\n# http://robotfindskitten.org/aw.cgi?main=software.rfk#rfksay\n#\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n\n     |\\\\_/|\n     |o o|__\n     --*--__\\\\\n     C_C_(___)\nEOC\n";

	var kitty = "##\n## A kitten of sorts, I think\n##\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n       (\"`-'  '-/\") .___..--' ' \"`-._\n         ` $eye_ $eye  )    `-.   (      ) .`-.__. `)\n         (_Y_.) ' ._   )   `._` ;  `` -. .-'\n      _.. `--'_..-_/   /--' _ .' ,4\n   ( i l ),-''  ( l i),'  ( ( ! .-'    \nEOC\n";

	var knight = "# Knight (Chess piece)\n#\n# from http://www.chessvariants.org/d.pieces/ascii.html\n#   by David Moeser\n#\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n  __/\"\"\"\\\\\n ]___ 0  }\n     /   }\n   /~    }\n   \\\\____/\n   /____\\\\\n  (______)\nEOC\n";

	var koala = "##\n## From the canonical koala collection\n##\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n       ___  \n     {~$eye_$eye~}\n      ( Y )\n     ()~*~()   \n     (_)-(_)   \nEOC\n";

	var kosh = "##\n## It's a Kosh Cow!\n##\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n      $thoughts\n  ___       _____     ___\n /   \\\\     /    /|   /   \\\\\n|     |   /    / |  |     |\n|     |  /____/  |  |     |     \n|     |  |    |  |  |     |\n|     |  | {} | /   |     |\n|     |  |____|/    |     |\n|     |    |==|     |     |\n|      \\\\___________/      |\n|                         |\n|                         |\nEOC\n";

	var lamb2 = "$the_cow = <<EOC;\n $thoughts\n  $thoughts\n  ,-''''-.\n (.  ,.   L        ___...__\n /$eye} ,-`  `'-==''``        ''._\n//{                           '`.\n\\\\_,X ,                         : )\n $tongue 7                          ;`\n    :                  ,       /\n     \\\\_,                \\\\     ;\n       Y   L_    __..--':`.    L\n       |  /| ````       ;  y  J\n       [ j J            / / L ;\n       | |Y \\\\          /_J  | |\n       L_J/_)         /_)   L_J\n      /_)               sk /_)\nEOC\n";

	var lamb = "$the_cow = <<EOC;\n                 $thoughts\n                  $thoughts  _,._\n                 __.'   _)\n                <_,)'.-\"$eye\\\\\n                  /' (    \\\\\n      _.-----..,-'   (`\"--^\n     //              |   $tongue\n    (|   `;      ,   |  \n      \\\\   ;.----/  ,/ \n       ) // /   | |\\\\ \\\\\n       \\\\ \\\\\\\\`\\\\   | |/ /\n        \\\\ \\\\\\\\ \\\\  | |\\\\/\nEOC\n";

	var lightbulb = "# lightbulb\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n$thoughts\n $thoughts\n         ,=;%\\$%%\\$X%%%%;/%%%%;=,\n     ,/\\$\\$+:-                -:+\\$\\$/,\n   :X\\$=                          =\\$X:\n ;M%.                              .%M;\n+#/                                  /#+\n\\##                                    M#\nH#,                     =;+/;,       ,#X\n.HM-       :\\@X+%H:   .%M%- .M#.     -M\\@.\n  /#%.     \\@#-  ,H\\@--MH, .;\\@\\$-    .%#+\n   .\\$M;    .+\\@X;, MM#\\@:/\\$X;.     ;M\\$,\n     =\\@H,     ,:+%H#M%;-       ,H\\@=\n      .\\$#;        -#H         =#\\$\n        %#;        \\#M        ;#%\n         H#-       \\##       -#H\n         ;#+       \\##       +#;\n          ;H+;;;;;;HH;;;;;;+H/\n           =H#\\@HHHHHHHHHH\\@#H=\n           =\\@#H%%%%%%%\\$HH\\@#\\@=\n           =\\@#X%%%%%%%\\$M###\\@=\n               =+%XHHX%+=\nEOC\n";

	var lobster = "# Lobster\n#   lobster jgs   10/96\n#   http://ascii.co.uk/art/lobster\n$the_cow = <<EOC;\n             $thoughts\n              $thoughts\n                             ,.---._\n                   ,,,,     /       `,\n                    \\\\\\\\\\\\\\\\   /    '\\\\_  ;\n                     |||| /\\\\/``-.__\\\\;'\n                     ::::/\\\\/_\n     {{`-.__.-'(`(^^(^^^(^ 9 `.========='\n    {{{{{{ { ( ( (  (   (-----:=\n     {{.-'~~'-.(,(,,(,,,(__6_.'=========.\n                     ::::\\\\/\\\\\n                     |||| \\\\/\\\\  ,-'/,\n                    ////   \\\\ `` _/ ;\n                   ''''     \\\\  `  .'\n                             `---'\nEOC\n";

	var lollerskates = "# LOLLERSKATES\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n        /\\\\O\n         /\\\\/\n        /\\\\\n       /  \\\\\n      LOL LOL\n:-D LOLLERSKATES :-D\nEOC\n";

	var lukeKoala = "##\n## From the canonical koala collection\n##\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts          .\n       ___   //\n     {~$eye_$eye~}// \n      ( Y )K/  \n     ()~*~()   \n     (_)-(_)   \n     Luke    \n     Skywalker\n     koala   \nEOC\n";

	var mailchimp = "# MailChimp\n#\n# view-source:http://mailchimp.com/\n$the_cow = <<EOC;\n$thoughts\n $thoughts\n    ______\n   / ___M ]__\nC{ ( o o )}\n    {     ••\n      \\\\___\n      ----´\nEOC\n";

	var mazeRunner = "# maze-runner.cow\n#\n#   a guy running through an ASCII maze\n#   found at http://pip.readthedocs.org/en/user_builds/pip/rtd-builds/latest/installing/\n#\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n      $thoughts\n       \\\\\n        \\\\\n         \\\\\n    \\\\     \\\\                     /\n     \\\\     \\\\                   /\n      \\\\     \\\\                 /\n       ]     \\\\               [    ,'|\n       ]      \\\\              [   /  |\n       ]___               ___[ ,'   |\n       ]  ]\\\\             /[  [ |:   |\n       ]  ] \\\\           / [  [ |:   |\n       ]  ]  ]         [  [  [ |:   |\n       ]  ]  ]__     __[  [  [ |:   |\n       ]  ]  ] ]\\\\ _ /[ [  [  [ |:   |\n       ]  ]  ] ] (#) [ [  [  [ :===='\n       ]  ]  ]_].nHn.[_[  [  [\n       ]  ]  ]  HHHHH. [  [  [\n       ]  ] /   `HH(\"N  \\\\ [  [\n       ]__]/     HHH  \"  \\\\[__[\n       ]         NNN         [\n       ]         N/\"         [\n       ]         N H         [\n      /          N            \\\\\n     /           q,            \\\\\n    /                           \\\\\nEOC\n\n";

	var mechAndCow = "$the_cow = <<EOC;\n      $thoughts                            |     |\n       $thoughts                        ,--|     |-.\n                         __,----|  |     | |\n                       ,;::     |  `_____' |\n                       `._______|    i^i   |\n                                `----| |---'| .\n                           ,-------._| |== ||//\n                           |       |_|P`.  /'/\n                           `-------' 'Y Y/'/'\n                                     .==\\ /_\\\n   ^__^                             /   /'|  `i\n   ($eyes)\\_______                   /'   /  |   |\n   (__)\\       )\\/\\             /'    /   |   `i\n    $tongue ||----w |           ___,;`----'.___L_,-'`\\__\n       ||     ||          i_____;----\\.____i\"\"\\____\\\nEOC\n";

	var meow = "##\n## A meowing tiger?\n##\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts ,   _ ___.--'''`--''//-,-_--_.\n      \\\\`\"' ` || \\\\\\\\ \\\\ \\\\\\\\/ / // / ,-\\\\\\\\`,_\n     /'`  \\\\ \\\\ || Y  | \\\\|/ / // / - |__ `-,\n    /\\$eye\"\\\\  ` \\\\ `\\\\ |  | ||/ // | \\\\/  \\\\  `-._`-,_.,\n   /  _.-. `.-\\\\,___/\\\\ _/|_/_\\\\_\\\\/|_/ |     `-._._)\n   `-'``/  /  |  // \\\\__/\\\\__  /  \\\\__/ \\\\\n    $tongue  `-'  /-\\\\/  | -|   \\\\__ \\\\   |-' |\n          __/\\\\ / _/ \\\\/ __,-'   ) ,' _|'\n         (((__/(((_.' ((___..-'((__,'\nEOC\n";

	var milk = "##\n## Milk from Milk and Cheese\n##\n$the_cow = <<EOC;\n $thoughts     ____________ \n  $thoughts    |__________|\n      /           /\\\\\n     /           /  \\\\\n    /___________/___/|\n    |          |     |\n    |  ==\\\\ /== |     |\n    |   $eye   $eye  | \\\\ \\\\ |\n    |     <    |  \\\\ \\\\|\n   /|          |   \\\\ \\\\\n  / |  \\\\_____/ |   / /\n / /|    $tongue    |  / /|\n/||\\\\|          | /||\\\\/\n    -------------|   \n        | |    | | \n       <__/    \\\\__>\nEOC\n";

	var minotaur = "$the_cow = <<\"EOC\";\n        $thoughts   ^__^\n         $thoughts  ($eyes)\n            (__)\n           /-||-\\\\\n           \\\\|\\\\/|/\n            o==o \n            ||||\n            ()()\nEOC\n";

	var monaLisa = "# Mona Lisa\n#\n# from http://www.heartnsoul.com/ascii_art/mona_lisa_ascii.htm\n$the_cow = <<EOC;\n          $thoughts\n           $thoughts\n\n!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!>''''''<!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n!!!!!!!!!!!!!!!!!!!!!!!!!!!!'''''`             ``'!!!!!!!!!!!!!!!!!!!!!!!!\n!!!!!!!!!!!!!!!!!!!!!!!!''`          .....         `'!!!!!!!!!!!!!!!!!!!!!\n!!!!!!!!!!!!!!!!!!!!!'`      .      :::::'            `'!!!!!!!!!!!!!!!!!!\n!!!!!!!!!!!!!!!!!!!'     .   '     .::::'                `!!!!!!!!!!!!!!!!\n!!!!!!!!!!!!!!!!!'      :          `````                   `!!!!!!!!!!!!!!\n!!!!!!!!!!!!!!!!        .,cchcccccc,,.                       `!!!!!!!!!!!!\n!!!!!!!!!!!!!!!     .-\"?\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$c,                      `!!!!!!!!!!!\n!!!!!!!!!!!!!!    ,ccc\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$,                     `!!!!!!!!!!\n!!!!!!!!!!!!!    z\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$;.                    `!!!!!!!!!\n!!!!!!!!!!!!    <\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$:.                    `!!!!!!!!\n!!!!!!!!!!!     \\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$h;:.                   !!!!!!!!\n!!!!!!!!!!'     \\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$h;.                   !!!!!!!\n!!!!!!!!!'     <\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$                   !!!!!!!\n!!!!!!!!'      `\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$F                   `!!!!!!\n!!!!!!!!        c\\$\\$\\$\\$???\\$\\$\\$\\$\\$\\$\\$P\"\"  \"\"\"??????\"                      !!!!!!\n!!!!!!!         `\"\" .,.. \"\\$\\$\\$\\$F    .,zcr                            !!!!!!\n!!!!!!!         .  dL    .?\\$\\$\\$   .,cc,      .,z\\$h.                  !!!!!!\n!!!!!!!!        <. \\$\\$c= <\\$d\\$\\$\\$   <\\$\\$\\$\\$=-=+\"\\$\\$\\$\\$\\$\\$\\$                  !!!!!!\n!!!!!!!         d\\$\\$\\$hcccd\\$\\$\\$\\$\\$   d\\$\\$\\$hcccd\\$\\$\\$\\$\\$\\$\\$F                  `!!!!!\n!!!!!!         ,\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$h d\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$                   `!!!!!\n!!!!!          `\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$<\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$'                    !!!!!\n!!!!!          `\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\"\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$P>                     !!!!!\n!!!!!           ?\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$??\\$c`\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$?>'                     `!!!!\n!!!!!           `?\\$\\$\\$\\$\\$\\$I7?\"\"    ,\\$\\$\\$\\$\\$\\$\\$\\$\\$?>>'                       !!!!\n!!!!!.           <<?\\$\\$\\$\\$\\$\\$c.    ,d\\$\\$?\\$\\$\\$\\$\\$F>>''                       `!!!\n!!!!!!            <i?\\$P\"??\\$\\$r--\"?\"\"  ,\\$\\$\\$\\$h;>''                       `!!!\n!!!!!!             \\$\\$\\$hccccccccc= cc\\$\\$\\$\\$\\$\\$\\$>>'                         !!!\n!!!!!              `?\\$\\$\\$\\$\\$\\$F\"\"\"\"  `\"\\$\\$\\$\\$\\$>>>''                         `!!\n!!!!!                \"?\\$\\$\\$\\$\\$cccccc\\$\\$\\$\\$??>>>>'                           !!\n!!!!>                  \"\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$F>>>>''                            `!\n!!!!!                    \"\\$\\$\\$\\$\\$\\$\\$\\$???>'''                                !\n!!!!!>                     `\"\"\"\"\"                                        `\n!!!!!!;                       .                                          `\n!!!!!!!                       ?h.\n!!!!!!!!                       \\$\\$c,\n!!!!!!!!>                      ?\\$\\$\\$h.              .,c\n!!!!!!!!!                       \\$\\$\\$\\$\\$\\$\\$\\$\\$hc,.,,cc\\$\\$\\$\\$\\$\n!!!!!!!!!                  .,zcc\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\n!!!!!!!!!               .z\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\n!!!!!!!!!             ,d\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$          .\n!!!!!!!!!           ,d\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$         !!\n!!!!!!!!!         ,d\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$        ,!'\n!!!!!!!!>        c\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$.       !'\n!!!!!!''       ,d\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$>       '\n!!!''         z\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$>\n!'           ,\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$>             ..\n            z\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$'           ;!!!!''`\n            \\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$F       ,;;!'`'  .''\n           <\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$>    ,;'`'  ,;\n           `\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$F   -'   ,;!!'\n            \"?\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$?\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$F     .<!!!'''       <!\n         !>    \"\"??\\$\\$\\$?C3\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\"\"     ;!'''          !!!\n       ;!!!!;,      `\"''\"\"????\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\\$\"\"   ,;-''               ',!\n      ;!!!!<!!!; .                `\"\"\"\"\"\"\"\"\"\"\"    `'                  ' '\n      !!!! ;!!! ;!!!!>;,;, ..                  ' .                   '  '\n     !!' ,;!!! ;'`!!!!!!!!;!!!!!;  .        >' .''                 ;\n    !!' ;!!'!';! !! !!!!!!!!!!!!!  '         -'\n   <!!  !! `!;! `!' !!!!!!!!!!<!       .\n   `!  ;!  ;!!! <' <!!!! `!!! <       /\n  `;   !>  <!! ;'  !!!!'  !!';!     ;'\n   !   !   !!! !   `!!!  ;!! !      '  '\n  ;   `!  `!! ,'    !'   ;!'\n      '   /`! !    <     !! <      '\n           / ;!        >;! ;>\n             !'       ; !! '\n          ' ;!        > ! '\n\nEOC\n";

	var moofasa = "##\n## MOOfasa.\n##\n$the_cow = <<EOC;\n       $thoughts    ____\n        $thoughts  /    \\\\\n          | ^__^ |\n          | ($eyes) |______\n          | (__) |      )\\\\/\\\\\n           \\\\____/|----w |\n                ||     ||\n\n\t         Moofasa\nEOC\n";

	var mooghidjirah = "$the_cow = <<EOC;\n $thoughts       $thoughts      $thoughts      \n  $thoughts        ^__^  $thoughts        \n    ^__^   ($eyes)   ^__^  \n    ($eyes)   (__)   ($eyes)   \n    (__)    $tongue    (__)   \noyo/:$tongue            $tongue:/oy+\n/mmmmm+   syyyyo  `ommmmm/\n smmmmms. -ymmy. .smmmmmo \n `+dmmmmd+``::``+dmmmmd+  \n   -ymmmmmh/``+hmmmmmy-   \n    `/hmmmmmhhmmmmmh/`    \n      `/hmmmmmmmmh/`      \n        `/hmmmmmd/        \n      `oo.`/dmmmmdo`      \n     `ymmd+``ommmmmy`     \n     smmmmd-  /mmmmms     \n    -mmmmm+    ommmmm-    \n    -ooooo`    .ooooo.     \nEOC\n";

	var moojira = "$the_cow = <<EOC;\n     $thoughts              \n      $thoughts    /ss/           \n   `oys:  .dmmd`  :syo`   \n   /dmmy   .//.   hmmd:   \n    -/:`          `:/-    \noyo/:.     ^__^     .:/oy+\n/mmmmm+   <($eyes\\)>  `ommmmm/\n smmmmms. -(__). .smmmmmo \n `+dmmmmd+``$tongue``+dmmmmd+  \n   -ymmmmmh/``+hmmmmmy-   \n    `/hmmmmmhhmmmmmh/`    \n      `/hmmmmmmmmh/`      \n        `/hmmmmmd/        \n      `oo.`/dmmmmdo`      \n     `ymmd+`VVmmmmmy`     \n     smmmmd-  /mmmmms     \n    -mmmmm+    ommmmm-    \n    -ooooo`    .ooooo.    \nEOC\n";

	var moose = "$the_cow = <<EOC;\n  $thoughts\n   $thoughts   \\\\_\\\\_    _/_/\n    $thoughts      \\\\__/\n           ($eyes)\\\\_______\n           (__)\\\\       )\\\\/\\\\\n            $tongue ||----- |\n               ||     ||\nEOC\n";

	var mule = "# Mule\n#\n# based on mule from http://rossmason.blogspot.com/2008/10/friday-ascii-art.html \n#\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts \n  /\\\\          /\\\\                               \n ( \\\\\\\\        // )                              \n  \\\\ \\\\\\\\      // /                               \n   \\\\_\\\\\\\\||||//_/                                \n     / _  _ \\\\/                                 \n                                               \n     |(o)(o)|\\\\/                                \n     |      | \\\\/                               \n     \\\\      /  \\\\/_____________________         \n      |____|     \\\\\\\\                  \\\\\\\\        \n     /      \\\\     ||                  \\\\\\\\       \n     \\\\ 0  0 /     |/                  |\\\\\\\\      \n      \\\\____/ \\\\    V           (       / \\\\\\\\     \n       / \\\\    \\\\     )          \\\\     /   \\\\\\\\    \n      / | \\\\    \\\\_|  |___________\\\\   /     \"\" \n                  ||  |     \\\\   /\\\\  \\\\          \n                  ||  /      \\\\  \\\\ \\\\  \\\\         \n                  || |        | |  | |         \n                  || |        | |  | |         \n                  ||_|        |_|  |_|         \n                 //_/        /_/  /_/          \nEOC\n";

	var mutilated = "##\n## A mutilated cow, from aspolito@csua.berkeley.edu\n##\n$the_cow = <<EOC;\n       $thoughts   \\\\_______\n v__v   $thoughts  \\\\   O   )\n ($eyes)      ||----w |\n (__)      ||     ||  \\\\/\\\\\n  $tongue\nEOC\n";

	var nyan = "# Nyan Cat\n#\n# from http://www.reddit.com/r/commandline/comments/2lb5ij/what_is_your_favorite_ascii_art/clt4ybl\n#\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n\n+      o     +              o   \n    +             o     +       +\no          +\n    o  +           +        +\n+        o     o       +        o\n-_-_-_-_-_-_-_,------,      o \n_-_-_-_-_-_-_-|   /\\\\_/\\\\  \n-_-_-_-_-_-_-~|__( ^ .^)  +     +  \n_-_-_-_-_-_-_-''  ''      \n+      o         o   +       o\n    +         +\no        o         o      o     +\n    o           +\n+      +     o        o      +    \nEOC\n\n";

	var octopus = "# octopus\n#   http://www.ascii-art.de/ascii/mno/octopus.txt\n$the_cow = <<EOC;\n        $thoughts               ___\n         $thoughts           .-'   `'.\n                    /         \\\\\n                    |         ;\n                    |         |           ___.--,\n           _.._     |0) ~ (0) |    _.---'`__.-( (_.\n    __.--'`_.. '.__.\\\\    '--. \\\\_.-' ,.--'`     `\"\"`\n   ( ,.--'`   ',__ /./;   ;, '.__.'`    __\n   _`) )  .---.__.' / |   |\\\\   \\\\__..--\"\"  \"\"\"--.,_\n  `---' .'.''-._.-'`_./  /\\\\ '.  \\\\ _.-~~~````~~~-._`-.__.'\n        | |  .' _.-' |  |  \\\\  \\\\  '.               `~---`\n         \\\\ \\\\/ .'     \\\\  \\\\   '. '-._)\n          \\\\/ /        \\\\  \\\\    `=.__`~-.\n     jgs  / /\\\\         `) )    / / `\"\".`\\\\\n    , _.-'.'\\\\ \\\\        / /    ( (     / /\n     `--~`   ) )    .-'.'      '.'.  | (\n            (/`    ( (`          ) )  '-;\n             `      '-;         (-'\nEOC\n";

	var okazu = "#\n# おかず\n#\n\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts                _, _ ,､\n    $thoughts          , - ´      `--、\n             ノ               丶\n           ／                  `､_\n         ,´                        、\n        ,'                          丶\n       ﾉ                             ヽ\n    ＿;＿＿＿＿＿＿＿＿＿＿＿＿＿＿＿',＿\n    ヽ三三三三三三三三三三三三三三三三三ﾉ\n      ヽ                              /\n       ヽ三三三三三三三三三三三三三三/\n         ＼                        ／\n           ＼三三三三三三三三三三／\n             `＜              ＞´\n               ｀丁三三三三丁´\n     ＿          ｀ ｰ----‐ ´\n  ／::/＿＿＿＿＿＿＿＿＿＿＿＿＿＿＿＿＿_\n（;;;ﾌ ｰ─----＝＝ === ニニニ 二二二三三三｣\n\n         ＿|＿ ＼  ＿ｌ＿＼  _＿|＿_ヽヽ\n          _|＿       ｜ヽ     __|\n        ／ |  ヽ     ﾉ  │   (__|\n        ＼ノ  ノ    ﾉ ヽﾉ     _ノ \nEOC\n";

	var owl = "##\n## An owl\n##\n$the_cow = <<EOC;\n         $thoughts\n          $thoughts\n           ___\n          (o o)\n         (  V  )\n        /--m-m-\nEOC\n";

	var pawn = "# Pawn (Chess piece)\n#\n# from http://www.chessvariants.org/d.pieces/ascii.html\n#   by David Moeser\n#\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n     __\n    (  )\n     ||\n    /__\\\\\n   (____)\nEOC\n";

	var periodicTable = "$the_cow = <<EOC;\n$thoughts\n $thoughts\n   1A   2A                                         3A  4A  5A  6A  7A  8A\n  -----                                                               -----\n1 | H |                                                               |He |\n  |---+----                                       --------------------+---|\n2 |Li |Be |                                       | B | C | N | O | F |Ne |\n  |---+---|                                       |---+---+---+---+---+---|\n3 |Na |Mg |3B  4B  5B  6B  7B |    8B     |1B  2B |Al |Si | P | S |Cl |Ar |\n  |---+---+---------------------------------------+---+---+---+---+---+---|\n4 | K |Ca |Sc |Ti | V |Cr |Mn |Fe |Co |Ni |Cu |Zn |Ga |Ge |As |Se |Br |Kr |\n  |---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---|\n5 |Rb |Sr | Y |Zr |Nb |Mo |Tc |Ru |Rh |Pd |Ag |Cd |In |Sn |Sb |Te | I |Xe |\n  |---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---|\n6 |Cs |Ba |Lu |Hf |Ta | W |Re |Os |Ir |Pt |Au |Hg |Tl |Pb |Bi |Po |At |Rn |\n  |---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---|\n7 |Fr |Ra |Lr |Rf |Db |Sg |Bh |Hs |Mt |Ds |Rg |Cn |Nh |Fl |Mc |Lv |Ts |Og |\n  -------------------------------------------------------------------------\n              -------------------------------------------------------------\n   Lanthanide |La |Ce |Pr |Nd |Pm |Sm |Eu |Gd |Tb |Dy |Ho |Er |Tm |Yb |Lu |\n              |---+---+---+---+---+---+---+---+---+---+---+---+---+---+---|\n   Actinide   |Ac |Th |Pa | U |Np |Pu |Am |Cm |Bk |Cf |Es |Fm |Md |No |Lr |\n              -------------------------------------------------------------\nEOC\n";

	var personalitySphere = "# Personality Sphere from Portal/Portal 2\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n      .-+\\$H###MM\\@MMMMM##\\@\\$+-,. ....\n-\\@\\$+%\\$+%HX+--..  .  . .,:X\\$/+/++\\$#:\n-#MXH\\$=                      \\$HXH#:\n .--,:#+   ,+\\$HMX =\\@\\@X%, . .X#:,,,\n     =#\\@\\$H :####H =####;,M%\\$#X\n     X###\\$ \\$####X =####H %###X\n    ;###X /###\\@\\$: ,+HM##H.+###;\n   :###;,X##%=;%H\\@H\\$;-;M#\\@-;###/\n  ,M##;.\\@##;-H#######M=.M##-:###-\n  ;##M ;##X \\@###H-=\\@###.;##X H##;\n  ;##M./##X.\\@###H:/M###-=##X X##;\n  -###;,M##:,\\@########+-H##; \\@##-\n   %##M==\\@##%==%HMH%::/M##+.X##+\n    %###/./###X+: -+\\$M##M=,X##+\n     X###X X####H +#####% \\@##H\n     :###H %####H +#####; X##;\n     /#\\$.  -HM##H /###\\@+.  +#\\$. .\n/HX%\\$X:      .,-, .-,.      =XX\\$H\\@-\n/#H+/+%+/+;=.          .=/%;;/;;+#+\n ..  .,-:XM#MM\\@\\@\\@\\@\\@\\@H\\@\\@M#\\@+=,.   ,,\nEOC\n";

	var pinballMachine = "# Pinball machine\n#\n# from http://ascii.co.uk/art/pinball\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n              /\\\\\n             <  \\\\\n             |\\\\  \\\\\n             | \\\\  \\\\\n             | .\\\\  >\n             |  .\\\\/|\n             |   .||\n             |    ||\n            / \\\\   ||\n           /,-.\\\\: ||\n          /,,  `\\\\ ||\n         /,  ', `\\\\||\n        /, *   ''/ |\n       /,    *,'/  |\n      /,     , /   |\n     / :    , /   .|\n    /\\\\ :   , /   /||\n   |\\\\ \\\\ .., /   / ||\n   |.\\\\ \\\\ . /   /  ||\n   |  \\\\ \\\\ /   /   ||\n   |   \\\\ /   /    |'\n   |\\\\o '|o  /\n   ||\\\\o |  /\n   || \\\\ | /\n   ||  \\\\|/\n   |'   ||\n        ||\n        ||\n        |'\nEOC\n";

	var psychiatrichelp2 = "use utf8;\n$the_cow = <<EOC;\n $thoughts      .------------------------.\n  $thoughts     |       PSYCHIATRIC      |\n   $thoughts    |         HELP  5¢       |\n    $thoughts   |________________________|\n     $thoughts  ||     .-\"\"\"--.         ||\n      $thoughts ||    /        \\\\.-.     ||\n        ||   |     ._,     \\\\    ||\n        ||   \\\\_/`-'   '-.,_/    ||\n        ||   (_   (' _)') \\\\     ||\n        ||   /|           |\\\\    ||\n        ||  | \\\\     __   / |    ||\n        ||   \\\\_).,_____,/}/     ||\n      __||____;_--'___'/ (      ||\n     |\\\\ ||   (__,\\\\\\\\    \\\\_/------||\n     ||\\\\||______________________||\n     ||||                        |\n     ||||       THE DOCTOR       |\n     \\\\|||         IS [IN]   _____|\n      \\\\||                  (______)\n jgs   `|___________________//||\\\\\\\\\n                           //=||=\\\\\\\\\n                           `  ``  `\nEOC\n";

	var psychiatrichelp = "$the_cow = <<EOC;\n        $thoughts         ____________________\n         $thoughts       |                    |\n          $thoughts      |     PSYCHIATRIC    |\n           $thoughts     |        HELP        |\n            $thoughts    |____________________|\n             $thoughts   ||  ,-..'``.        ||\n              $thoughts  || (,-..'`. )       ||\n                 ||   )-c - `)\\\\      ||\n   ,.,._.-.,_,.,-||,.(`.--  ,`',.-,_,||.-.,.,-,._.\n              ___||____,`,'--._______||\n             |`._||______`'__________||\n             |   ||     __           ||\n             |   ||    |.-' ,|-      ||\n   _,_,,..-,_|   ||    ._)) `|-      ||,.,_,_.-.,_\n            . `._||__________________||   ____    .\n     .              .           .     . <.____`>\n   .SSt  .      .     .      .    .   _.()`'()`'  .\nEOC\n   ";

	var pterodactyl = "# pterodactyl.cow\n#\n#   a pterodactyl with its mouth open\n#\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n      $thoughts\n                                                                                 -/- \n                                                                              -/ --/    \n                                                                            /- -  /     \n                                                                         //      /      \n                                                                        /       /       \n                                                                      //       /        \n                                                                    //        /         \n                                                                  //          /         \n                                                                ///           /         \n                                                               //            /          \n                                                              //            /           \n                                                             //          . ./           \n                                                             //       .    /            \n                                                             //    .      /             \n                                                             //  .       /              \n                                                            // .         /              \n                                                          (=>            /              \n                                                         (==>            /              \n                                                          (=>            /              \n             -_                                           //.           /               \n             \\\\\\\\-_                                        //   .         /               \n              \\\\ \\\\_-_                                     //     .       /               \n               \\\\_ \\\\_--_                                 //        . . . /               \n                 \\\\_ \\\\_ -_                              //              /                \n                   \\\\_ \\\\_ (O)-___                      //               /                \n                     \\\\ _\\\\   __  --__                  /                /                \n                     _/    \\\\  ----__--____          //                 /                \n                   _/  _/   \\\\       -------       //                  /                 \n                 _/ __/ \\\\\\\\   \\\\\\\\                  /                   /                  \n               _/ _/      \\\\\\\\   \\\\\\\\              //                   /                   \n              -__/          \\\\\\\\   \\\\\\\\\\\\          //                   /                    \n                              \\\\\\\\    \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\//   -                /                    \n                                \\\\\\\\         _/         -            /                    \n                                  \\\\\\\\                      -        \\\\                    \n                                    \\\\\\\\\\\\                       -     \\\\                   \n                                        \\\\\\\\                       -   \\\\                  \n                                          \\\\\\\\\\\\                         \\\\--__             \n                                           | \\\\\\\\                            \\\\__________  \n                                            |  \\\\\\\\\\\\\\\\                ___      _________-\\\\\\\\\n                                            |    \\\\\\\\\\\\\\\\\\\\                \\\\--__/____        \n                                            |        \\\\\\\\\\\\\\\\________---\\\\-    ______-----   \n                                             |                   /    \\\\--  \\\\_______     \n                                             |                   /       \\\\-_________\\\\   \n                                             \\\\                   /                  \\\\\\\\  \n                                             \\\\                 ./                       \n                                             \\\\            .     /                       \n                                              \\\\        .       /                        \n                                              \\\\    .           //                       \n                                              \\\\                /                        \n                                              |__              /                        \n                                              \\\\==              /                        \n                                               \\\\\\\\              \\\\                        \n                                                \\\\\\\\  .          \\\\                        \n                                                  \\\\\\\\    .  .   \\\\                        \n                                                   \\\\           .\\\\                       \n                                                   \\\\\\\\            \\\\                      \n                                                     \\\\           \\\\                      \n                                                      \\\\\\\\          \\\\                     \n                                                        \\\\\\\\         \\\\                    \n                                                          \\\\         \\\\--                 \n                                                           \\\\\\\\          \\\\                \n                                                             \\\\\\\\         \\\\\\\\\\\\\\\\            \n                                                               \\\\\\\\\\\\\\\\_________\\\\\\\\\\\\         \nEOC\n\n";

	var queen = "# Queen (Chess piece)\n#\n# from http://www.chessvariants.org/d.pieces/ascii.html\n#   by David Moeser\n#\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n     ()\n   <~~~~>\n    \\\\__/\n   (____)\n    |  |\n    |  |\n    |__|\n   /____\\\\\n  (______)\n (________)\nEOC\n";

	var R2D2 = "# R2-D2\n#\n# from http://www.ascii-art.de/ascii/s/starwars.txt\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n         _____\n       .\\'/L|__`.\n      / =[_]O|` \\\\\n      |\\\"+_____\\\":|\n    __:='|____`-:__\n   ||[] ||====| []||\n   ||[] | |=| | []||\n   |:||_|=|U| |_||:|\n   |:|||]_=_ =[_||:| LS\n   | |||] [_][]C|| |\n   | ||-\\'\\\"\\\"\\\"\\\"\\\"`-|| |\n   /|\\\\\\\\_\\\\_|_|_/_//|\\\\\n  |___|   /|\\\\   |___|\n  `---\\'  |___|  `---\\'\n         `---'\nEOC\n";

	var radio = "# radio from Portal\n# via http://pastebin.com/1AZwKrKp\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n                    ;=\n                    /=\n                    ;=\n                    /=\n                    ;=\n                    /=\n                    ;=\n                    /=\n             ,--==-:\\$;\n         ,/\\$@#######\\@X+-\n      ./@###############X=\n     /M#####X+/;;;;+H#####\\$.\n    %####M/;+H\\@XX@@%;;\\@####\\@,\n   +####H=+##\\$,--,=M#X-%####\\@.\n  -####X,X\\@HHXH##MXHXXH-+####\\$\n  X###\\@.X/\\$M\\$:####\\$=\\@X/X,X####-\n .####:+\\$:##\\@:####\\$:##H/X=####%\n -%%\\$%,+==%\\$+-\\$+:\\$;-\\$\\$%-+,/\\$%%+\n -/+%%X\\$XX\\$\\$\\$\\$\\$\\$\\$%\\$\\$\\$%\\$X\\$X\\$%+/-\nEOC\n";

	var ren = "##\n## Ren \n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n    ____  \n   /# /_\\\\_\n  |  |/$eye\\\\$eye\\\\\n  |  \\\\\\\\_/_/\n / |_   |  \n|  ||\\\\_ ~| \n|  ||| \\\\/  \n|  |||_    \n \\\\//  |    \n  ||  |    \n  ||_  \\\\   \n  \\\\_|  o|  \n  /\\\\___/   \n /  ||||__ \n    (___)_)\nEOC\n";

	var renge = "##\n## Nyanpasu~\n##\n$the_cow = <<EOC;\n     $thoughts               _\n      $thoughts            ´   ＼   __\n       $thoughts        ／ ／⌒\\\\ | ／   ＼\n   f|{r、       | /     '|/ ／⌒＼＼\n   ||J |        \\\\/＞--＜\\\\/ /--    |\n(＼|`` し]ﾄ----／          ⌒` ＼| /\n ＼      ﾉ\\\\   /                ＼|/\\\\   --、___\n  ゛    /  ＼/      /     |         \\\\/_       ﾉ\n   \\\\、/\\\\_／/ｲ    ,/'|    /\\\\ 、        Ⅵ   __／\n    [\\\\/   \\\\/_|   /\\\\|/|   |-]  、     く-く\n    |      \\\\/|  |/___ﾉ\\\\  /\\\\___ \\\\     /   ＼\n    {/      <|小| _ﾒﾘ  \\\\/  _ﾒﾘ` \\\\   ｜|   |\n     \\\\        ｜| \\\\/ｿ      \\\\/ｿ  ﾉ / /\\\\|＼_/\n      \\\\       ｜|              /_ｲ\\\\/\n       \\\\      ｜|     /ヽ      / /ﾉ\n        \\\\     ｜/\\\\   └-     ,/ /'\n         \\\\    ｜ |／>> r -=≦{{/ /ﾆ=_\n          \\\\   人 | ／ｨ|     /ﾚ/__   ﾉﾆ-、\n           ＼   \\\\|/  Xﾉ    / /   入//⌒Yﾊ\n             \\\\  /し ｜`---' //  /  \\\\ﾆﾆﾆﾉ|\n              ＼/  / \\\\  --ｱ ｜  |   | _]|\n               ｜ /   \\\\/\\\\/  ｜  |   |___|\n               r勺    ｜_｜ ｜  |   |  ||\n               |`7    ｜ ｜ ｜  |   |   |\nEOC\n";

	var robot = "# Robot\n#\n# based on rfksay by Andrew Northern\n# http://robotfindskitten.org/aw.cgi?main=software.rfk#rfksay\n#\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n\n     [-]\n     (+)=C\n     | |\n     OOO\nEOC\n";

	var robotfindskitten = "# Robot finds kitten <3\n#\n# based on rfksay by Andrew Northern\n# http://robotfindskitten.org/aw.cgi?main=software.rfk#rfksay\n#\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n\n    [-]   |\\\\_/|\n    (+)=C |o o|__\n    | |   --*--__\\\\\n    OOO   C_C_(___)\nEOC\n";

	var roflcopter = "$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n ROFL:ROFL:ROFL:ROFL\n         _^___\n L    __/   $eyes \\\\    \nLOL===__        \\\\ \n L      \\\\________]\n         I   I  $tongue\n        --------/\nEOC\n";

	var rook = "# Rook (Chess piece)\n#\n# from http://www.chessvariants.org/d.pieces/ascii.html\n#   by David Moeser\n#\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n\n   WWWWWW\n    |  |\n    |  |\n    |__|\n   /____\\\\\n  (______)\nEOC\n";

	var sachiko = "#\n# プロデューサーさんは独特の変わったセンスをしてますね！\n#\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n       $thoughts\n             , -――- 、\n          ／          ヽ、\n        /爻ﾉﾘﾉﾊﾉﾘlﾉ ゝ  l\n     ＜ﾉﾘﾉ‐'    ｰ  ﾘ ＞ }\n        l ﾉ ┃    ┃ l ﾉ  ﾉ\n        l人   r‐┐   !ﾉ＾)\n           ゝ ` ´ ‐＜´\nEOC\n";

	var satanic = "##\n## Satanic cow, source unknown.\n##\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts  (__)  \n         (\\\\/)  \n  /-------\\\\/    \n / | 666 ||$tongue  \n*  ||----||      \n   ~~    ~~      \nEOC\n";

	var seahorseBig = "# large seahorse\n#\n# adapted from http://www.chris.com/ascii/index.php?art=animals/other%20(water)\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n                  ,\n         ___     /^\\\\   ,\n        `\\  \\'...`   \\\\_/^\\\\\n          ) ~     ',    /__,\n         /       ,.    ,, /___,\n        (  .-.   \\'.\\'. /// ___/\n         ) .-.\\'  .`.`///-.\\'.\n        / ( o )  .\\\"\\\". ====) \\\\\n       (   \\'-`   \\\\  |\\'~~~`  u\\\\,\n        \\\\ _~  .\\\"\\\"\\\"` |~|^u^ u^(\\\"\\\"\n        //  .\"     /~/^ u^ u^\\\n       // .\"      /~  u^ u  ^u\\      _\n      // .\"      /~/U^ U^ U^ ^(     / )\n     /` .\"       |~  U^ U^ ^ U^\\   /) _)\n   ./` .\"        |~|^ U^ ^U ^ U(  / _  _)\n  ;.`.\"          |~ ^U ^ U^ U ^/ /)_ =  _)\n   \\\"\\\"            |~|^ ^U ^ ^ U(_/_    )- _)\n                 |~ U ^ ^U ^U ^ )   =    _)\n                 \\\\~|^ U U^ U ^ =  ~ )  - _)\n                  \\\\ U ^U ^ ^U^_)     =  _)\n                   \\\",^U^ ^U ^/ \\\\)_~   -_)\n                     \\\".u^u ^|   \\\\_  = _)\n                      ).u ^u|    \\\\)  _)\n                      \\\\u ^u^(     \\\\__)\n                       )^u ^u\\\\\n                       \\\\u ^u ^|\n             ____       )^u ^u|\n          ,-`    '-.    )u ^u^|\n         /  .---. ' \\\\  / ^ u^/\n        |  ;  `  '  | /u^u ^/\n        |  ;  '-` . `:u^u^u/\n        \\\\.\\'^\\'._   _.`u ^.-`\n         \\\\_.~=_```-.^.-\\\"\n           \\'\\\"------\\\"`\n\nEOC\n";

	var seahorse = "# seahorse\n#\n# adapted from http://www.chris.com/ascii/index.php?art=animals/other%20(water)\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n\n      (\\\\(\\\\/\n  .-._)oo  '_\n  \\'---.     .\\'\\\\\n       )    \\\\.-\\'\\\\\n      /__ ;     (\n      |__ : /'._/\n       \\\\_  (\n       .,)  )\n       \\'-.-\\'\n\nEOC\n";

	var sheep = "##\n## The non-flaming sheep.\n##\n$the_cow = <<EOC\n  $thoughts\n   $thoughts\n       __     \n      U${eyes}U\\\\.'\\@\\@\\@\\@\\@\\@`.\n      \\\\__/(\\@\\@\\@\\@\\@\\@\\@\\@\\@\\@)\n        $tongue (\\@\\@\\@\\@\\@\\@\\@\\@)\n           `YY~~~~YY'\n            ||    ||\nEOC\n";

	var shikato = "$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n\n     Lｰ'{r ｧjｰノ\n      _`)-ﾑ{\n    /´::( ･)ヽ-- ､\n   {::::::::::::::}\n   ゝ:::::.ノー-\n     しｿ¨UU\nEOC\n";

	var shrug = "$the_cow = <<EOC;\n  $thoughts\n¯\\\\_(ツ)_/¯\nEOC\n";

	var skeleton = "##\n## This 'Scowleton' brought to you by one of \n## {appel,kube,rowe}@csua.berkeley.edu\n##\n$the_cow = <<EOC;\n          $thoughts      (__)      \n           $thoughts     /$eyes|  \n            $thoughts   (_\"_)*+++++++++*\n                   //I#\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\I\\\\\n                   I[I|I|||||I I `\n                   I`I'///'' I I\n                   I I       I I\n                   ~ ~       ~ ~\n                     Scowleton\nEOC\n";

	var small = "##\n## A small cow, artist unknown\n##\n$eyes = \"..\" unless ($eyes);\n$the_cow = <<EOC;\n       $thoughts   ,__,\n        $thoughts  ($eyes)____\n           (__)    )\\\\\n            $tongue||--|| *\nEOC\n";

	var smilingOctopus = "# \n$the_cow = <<EOC;\n      $thoughts\n       $thoughts\n        $thoughts                                     ,\n                                            ,o\n                                            :o\n                   _....._                  `:o\n                 .\\'       ``-.                \\\\o\n                /  _      _   \\\\                \\\\o\n               :  /*\\\\    /*\\\\   )                ;o\n               |  \\\\_/    \\\\_/   /                ;o\n               (       U      /                 ;o\n                \\\\  (\\\\_____/) /                  /o\n                 \\\\   \\\\_m_/  (                  /o\n                  \\\\         (                ,o:\n                  )          \\\\,           .o;o\\'           ,o\\'o\\'o.\n                ./          /\\\\o;o,,,,,;o;o;\\'\\'         _,-o,-\\'\\'\\'-o:o.\n .             ./o./)        \\\\    \\'o\\'o\\'o\\'\\'         _,-\\'o,o\\'         o\n o           ./o./ /       .o \\\\.              __,-o o,o\\'\n \\\\o.       ,/o /  /o/)     | o o\\'-..____,,-o\\'o o_o-\\'\n `o:o...-o,o-\\' ,o,/ |     \\\\   \\'o.o_o_o_o,o--\\'\\'\n .,  ``o-o\\'  ,.oo/   \\'o /\\\\.o`.\n `o`o-....o\\'o,-\\'   /o /   \\\\o \\\\.                       ,o..         o\n   ``o-o.o--      /o /      \\\\o.o--..          ,,,o-o\\'o.--o:o:o,,..:o\n                 (oo(          `--o.o`o---o\\'o\\'o,o,-\\'\\'\\'        o\\'o\\'o\n                  \\\\ o\\\\              ``-o-o\\'\\'\\'\\'\n   ,-o;o           \\\\o \\\\\n  /o/               )o )  Carl Pilcher\n (o(               /o /                |\n  \\\\o\\.       ...-o\\'o /              \\\\   |\n    \\\\o`o`-o\\'o o,o,--\\'       ~~~~~~~~\\\\~~|~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n      ```o--\\'\\'\\'                       \\\\| /\n                                       |/\n ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~|~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n                                       |\n ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n\nEOC\n";

	var snoopy = "##\n## acsii picture From: kwok@menpachi.nmfs.hawaii.edu (William Kwok)\n## from http://www.ascii-art.de/ascii/s/snoopy.txt\n$the_cow = <<EOC;\n $thoughts\n  $thoughts          , ----.\n   $thoughts        -  -     `\n      ,__.,'           \\\\\n    .'                 *`\n   /       $eye   $eye     / **\\\\\n  .                 / ****.\n  |    mm           | ****|\n   \\\\                | ****|\n    ` ._______      \\\\ ****/\n              \\\\      /`---'\n               \\\\___(\n               /~~~~\\\\\n              /      \\\\\n             /      | \\\\\n            |       |  \\\\\n  , ~~ .    |, ~~ . |  |\\\\\n ( |||| )   ( |||| )(,,,)`\n( |||||| )-( |||||| )    | ^\n( |||||| ) ( |||||| )    |'/\n( |||||| )-( |||||| )___,'-\n ( |||| )   ( |||| )\n  ` ~~ '     ` ~~ '\nEOC\n";

	var snoopyhouse = "##\n## acsii picture from http://www.ascii-art.de/ascii/s/snoopy.txt\n##\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts       __---__                         ______\n         $thoughts     /    ___\\\\_             o  O  O _(      )__\n              /====(_____\\\\___---_  o        _(           )_\n             |                    \\\\        (_  AI-YA!!!!   )\n             |                     |@        (_  Shot      _)\n              \\\\       ___         /           (__  Again!__)\n \\\\ __----____--_\\\\____(____\\\\_____/                (______)\n==|__----____--______|\n /              /    \\\\____/)_\n              /        ______)\n             /           |  |\n            |           _|  |\n       ______\\\\______________|______\n      /                    *   *   \\\\\n     /_____________*____*___________\\\\\n     /   *     *                    \\\\\n    /________________________________\\\\\n    / *                              \\\\\n   /__________________________________\\\\\n        |                        |\n        |________________________|\n        |                        |\n        |________________________|\nEOC\n";

	var snoopysleep = "##\n## picture from http://www.ascii-art.de/ascii/ab/beagle.txt\n## \n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n   $thoughts     O_      __)(\n       ,'  `.   (_\".`.\n      :      :    /|`\n      |      |   ((|_  ,-.\n      ; -   /:  ,'  `:(( -\\\\\n     /    -'  `: ____ \\\\\\\\\\\\-:\n    _\\\\__   ____|___  \\\\____|_\n   ;    | |        '-`      :\n  :_____|:|__________________:\n  ;     |:|                  :\n :      |:|                   :\n ;_______`'___________________:\n:                              :\n|______________________________|\n `---.--------------------.---'\n     |____________________|\n     |                    |\n     |____________________|\n     |                    |\n   _\\\\|_\\\\|_\\\\/(__\\\\__)\\\\__\\\\//_|(_\nEOC\n";

	var spidercow = "$the_cow = <<EOC;\n          $thoughts     (\n           $thoughts     )\n            $thoughts   (\n         /\\\\  .-\"\"\"\"-.  /\\\\\n        //\\\\\\\\/  ,,,,  \\\\//\\\\\\\\\n        |/\\\\| ,;;;;;;, |/\\\\|\n        //\\\\\\\\\\\\;-\"\"\"\"-;///\\\\\\\\\n       //  \\\\/   ..   \\\\/  \\\\\\\\\n      (| ,-_| \\\\ || / |_-, |)\n        //`__(\\\\(__)/)__`\\\\\\\\\n       // /.-\\\\`($eyes)'/-.\\\\ \\\\\\\\\n      (\\\\ |)   ')  ('   (| /)\n       ` (|   (o  o)   |) `\n         \\\\)    `--'    (/\n                $tongue\nEOC\n";

	var squid = "#\n# これスプラトゥーン感あるね\n#\n\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts                                                           ＿＿＿ノ^l\n      $thoughts                                            ＿,,ノ``ｰ-'￣￣        ｌ\n                                                 く                       /\n                                                  `ヽ,   __､-'           /\n                                                    __＞‐´               |\n                                           ._,;‐''``              ,     /\n                                         _;\"                     /     /\n                                       ／                       /     く\n                                     ／                        /       |\n                                   ／                        ／       ｌ\n                                 ノ                        ／￣ヽ     /\n                                /                        ／     ） _ノ\n                            ,r'″ヽ、                   ／        ￣\n                           /      ヽ                 ／\n                        ＿ﾉ        `r            _､‐'\n                      ／          _l,_       _､‐'\n                 __,r'          ／r;;,ヽ   ／\n               ,/              ｜.;●,;;|  ノ\n              ノ ／  ／／       ヽ､!!!ﾞﾉ \"\n            ／ ／／／  ／／___,r''\"￣\n           / ／ / / /／ / /\n      ___／／/／／／ ／／/\n  ／￣＿_／／／/ / ／／／\n l ／´___／／／／／／ /\n しレ\"／／/ /  ／／//／\n      / ,/ / ／／／ /\n      ﾚ'   ﾚ'／／ ／\n           ／l｜l/\n          ｜|ﾚ'lノ\n           レ'\nEOC\n";

	var squirrel = "$the_cow = <<EOC;\n  $thoughts\n     $thoughts\n                  _ _\n       | \\__/|  .~    ~.\n       /$eyes `./      .'\n      {o__,   \\    {\n        / .  . )    \\\n        `-` '-' \\    }\n       .(   _(   )_.'\n      '---.~_ _ _|\n                                                     \nEOC\n";

	var stegosaurus = "##\n## A stegosaur with a top hat?\n##\n$the_cow = <<EOC;\n$thoughts                             .       .\n $thoughts                           / `.   .' \" \n  $thoughts                  .---.  <    > <    >  .---.\n   $thoughts                 |    \\\\  \\\\ - ~ ~ - /  /    |\n         _____          ..-~             ~-..-~\n        |     |   \\\\~~~\\\\.'                    `./~~~/\n       ---------   \\\\__/                        \\\\__/\n      .'  $eye    \\\\     /               /       \\\\  \" \n     (_____,    `._.'               |         }  \\\\/~~~/\n      `----.          /       }     |        /    \\\\__/\n            `-.      |       /      |       /      `. ,~~|\n                ~-.__|      /_ - ~ ^|      /- _      `..-‘ / \\\\  /\\\\\n                     |     /        |     /     ~-.     `-/ _ \\\\/__\\\\\n                     |_____|        |_____|         ~ - . _ _ _ _ _>\nEOC\n";

	var stimpy = "##\n## Stimpy!\n##\n$the_cow = <<EOC;\n  $thoughts     .    _  .    \n   $thoughts    |\\\\_|/__/|    \n       / / \\\\/ \\\\  \\\\  \n      /__|$eye||$eye|__ \\\\ \n     |/_ \\\\_/\\\\_/ _\\\\ |  \n     | | (____) | ||  \n     \\\\/\\\\___/\\\\__/  // \n     (_/         ||\n      |          ||\n      |          ||\\\\   \n       \\\\        //_/  \n        \\\\______//\n       __ || __||\n      (____(____)\nEOC\n";

	var sudowoodo = "# Sudowoodo (Pokémon)\n#\n# https://gist.github.com/rzabcik/9233650/\n#\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n     _              __\n    / `\\\\  (~._    ./  )\n    \\\\__/ __`-_\\\\__/ ./\n   _ \\\\ \\\\/  \\\\   \\\\ |_   __\n (   )  \\\\__/ -^    \\\\ /  \\\\\n  \\\\_/ \"  \\\\  | o  o  |.. /  __\n       \\\\. --' ====  /  || /  \\\\\n         \\\\   .  .  |---__.\\\\__/\n         /  :     /   |   |\n         /   :   /     \\\\_/\n      --/ ::    (\n     (  |     (  (____\n   .--  .. ----**.____)\n   \\\\___/\nEOC\n";

	var supermilker = "##\n## A cow being milked, probably from Lars Smith (lars@csua.berkeley.edu)\n##\n$the_cow = <<EOC;\n  $thoughts   ^__^\n   $thoughts  ($eyes)\\\\_______        ________\n      (__)\\\\       )\\\\/\\\\    |Super |\n       $tongue ||----W |       |Milker|\n          ||    UDDDDDDDDD|______|\nEOC\n";

	var surgery = "##\n## A cow operation, artist unknown\n##\n$the_cow = <<EOC;\n          $thoughts           \\\\  / \n           $thoughts           \\\\/  \n               (__)    /\\\\         \n               ($eyes)   O  O        \n               _\\\\/_   //         \n         *    (    ) //       \n          \\\\  (\\\\\\\\    //       \n           \\\\(  \\\\\\\\    )                              \n            (   \\\\\\\\   )   /\\\\                          \n  ___[\\\\______/^^^^^^^\\\\__/) o-)__                     \n |\\\\__[=======______//________)__\\\\                    \n \\\\|_______________//____________|                    \n     |||      || //||     |||\n     |||      || @.||     |||                        \n      ||      \\\\/  .\\\\/      ||                        \n                 . .                                 \n                '.'.`                                \n\n            COW-OPERATION                           \nEOC\n";

	var tableflip = "$the_cow = <<EOC;\n  $thoughts\n(╯°□°）╯︵ ┻━┻\nEOC\n";

	var taxi = "# Taxi cab\n#\n# from http://ascii.co.uk/art/taxi\n$the_cow = <<EOC;\n     $thoughts\n      $thoughts\n                   [\\\\\n              .----' `-----.\n             //^^^^;;^^^^^^`\\\\\n     _______//_____||_____()_\\\\________\n    /826    :      : ___              `\\\\\n   |>   ____;      ;  |/\\\\><|   ____   _<)\n  {____/    \\\\_________________/    \\\\____}\n       \\\\ '' /                 \\\\ '' /\n jgs    '--'                   '--'\nEOC\n";

	var telebears = "##\n## A cow performing an unnatural act, artist unknown.\n##\n$the_cow = <<EOC;\n      $thoughts                _\n       $thoughts              (_)   <-- TeleBEARS\n        $thoughts   ^__^       / \\\\\n         $thoughts  ($eyes)\\\\_____/_\\\\ \\\\\n            (__)\\\\  you  ) /\n             $tongue ||----w ((\n                ||     ||>> \nEOC\n";

	var template = "# \n$the_cow = <<EOC;\n$thoughts\n $thoughts\nEOC\n";

	var threader = "$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n         $thoughts\n             ＿＿＿＿\n           ／＿＿＿＿＼\n         ／／ (⌒ ⌒ ヽ＼＼\n        ｜｜  ﾉz(⌒ )| ｜｜\n        ｜｜ <   (_ノ ｜｜\n        ｜｜  L_／ )  ｜｜\n         ＼＼ /＿／  ／／\n           ＼⌒ )  (⌒ ／\n           ／／    ＼＼\n           ＼＼_  _／／\n             ﾍ＿)(＿/\n             ｜＝＝｜\n              ＼三／\n                ∧\n              ／  ＼\n              ＼  ／\n                Ｖ\nEOC\n";

	var threecubes = "# Three cubes\n#\n# from http://www.reddit.com/r/commandline/comments/2lb5ij/what_is_your_favorite_ascii_art/cltcqs1\n#   also available at https://gist.github.com/th3m4ri0/6e3f631866da31d05030\n# \n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n        ____________\n       /\\\\  ________ \\\\\n      /  \\\\ \\\\______/\\\\ \\\\\n     / /\\\\ \\\\ \\\\  / /\\\\ \\\\ \\\\\n    / / /\\\\ \\\\ \\\\/ / /\\\\ \\\\ \\\\\n   / / /__\\\\ \\\\ \\\\/_/__\\\\_\\\\ \\\\__________\n  / /_/____\\\\ \\\\__________  ________ \\\\\n  \\\\ \\\\ \\\\____/ / ________/\\\\ \\\\______/\\\\ \\\\\n   \\\\ \\\\ \\\\  / / /\\\\ \\\\  / /\\\\ \\\\ \\\\  / /\\\\ \\\\ \\\\\n    \\\\ \\\\ \\\\/ / /\\\\ \\\\ \\\\/ / /\\\\ \\\\ \\\\/ / /\\\\ \\\\ \\\\\n     \\\\ \\\\/ / /__\\\\_\\\\/ / /__\\\\ \\\\ \\\\/_/__\\\\_\\\\ \\\\\n      \\\\  /_/______\\\\/_/____\\\\ \\\\___________\\\\\n      /  \\\\ \\\\______/\\\\ \\\\____/ / ________  /\n     / /\\\\ \\\\ \\\\  / /\\\\ \\\\ \\\\  / / /\\\\ \\\\  / / /\n    / / /\\\\ \\\\ \\\\/ / /\\\\ \\\\ \\\\/ / /\\\\ \\\\ \\\\/ / /\n   / / /__\\\\ \\\\ \\\\/_/__\\\\_\\\\/ / /__\\\\_\\\\/ / /\n  / /_/____\\\\ \\\\_________\\\\/ /______\\\\/ /\n  \\\\ \\\\ \\\\____/ / ________  __________/\n   \\\\ \\\\ \\\\  / / /\\\\ \\\\  / / /\n    \\\\ \\\\ \\\\/ / /\\\\ \\\\ \\\\/ / /\n     \\\\ \\\\/ / /__\\\\_\\\\/ / /\n      \\\\  / /______\\\\/ /\n       \\\\/___________/\nEOC\n\n";

	var toaster = "# Toaster\n#   http://ascii.co.uk/art/toaster \n$the_cow = <<EOC;\n   $thoughts                     .___________.\n    $thoughts                    |           |\n     $thoughts    ___________.   |  |    /~\\\\ |\n         / __   __  /|   | _ _   |_| |\n        / /:/  /:/ / |   !________|__!\n       / /:/  /:/ /  |            |\n      / /:/  /:/ /   |____________!\n     / /:/  /:/ /    |\n    / /:/  /:/ /     |\n   /  ~~   ~~ /      |\n   |~~~~~~~~~~|      |\n   |    ::    |     /\n   |    ==    |    /\n   |    ::    |   /\n   |    ::    |  /\n   |    ::  @ | /\n   !__________!/\nEOC\n";

	var tortoise = "# Tortoise\n# from http://svn.haxx.se/tsvn/archive-2005-06/1030.shtml (accessed 9/11/2014)\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts       ___\n      oo  // \\\\\\\\\n     (_,\\\\/ \\\\_/ \\\\\n       \\\\ \\\\_/_\\\\_/>\n       /_/   \\\\_\\\\\nEOC\n";

	var turkey = "##\n## Turkey!\n##\n$the_cow = <<EOC;\n  $thoughts                                  ,+*^^*+___+++_\n   $thoughts                           ,*^^^^              )\n    $thoughts                       _+*                     ^**+_\n     $thoughts                    +^       _ _++*+_+++_,         )\n              _+^^*+_    (     ,+*^ ^          \\\\+_        )\n             {       )  (    ,(    ,_+--+--,      ^)      ^\\\\\n            { (\\@)    } f   ,(  ,+-^ __*_*_  ^^\\\\_   ^\\\\       )\n           {:;-/    (_+*-+^^^^^+*+*<_ _++_)_    )    )      /\n          ( /  (    (        ,___    ^*+_+* )   <    <      \\\\\n           U _/     )    *--<  ) ^\\\\-----++__)   )    )       )\n            (      )  _(^)^^))  )  )\\\\^^^^^))^*+/    /       /\n          (      /  (_))_^)) )  )  ))^^^^^))^^^)__/     +^^\n         (     ,/    (^))^))  )  ) ))^^^^^^^))^^)       _)\n          *+__+*       (_))^)  ) ) ))^^^^^^))^^^^^)____*^\n          \\\\             \\\\_)^)_)) ))^^^^^^^^^^))^^^^)\n           (_             ^\\\\__^^^^^^^^^^^^))^^^^^^^)\n             ^\\\\___            ^\\\\__^^^^^^))^^^^^^^^)\\\\\\\\\n                  ^^^^^\\\\uuu/^^\\\\uuu/^^^^\\\\^\\\\^\\\\^\\\\^\\\\^\\\\^\\\\^\\\\\n                     ___) >____) >___   ^\\\\_\\\\_\\\\_\\\\_\\\\_\\\\_\\\\)\n                    ^^^//\\\\\\\\_^^//\\\\\\\\_^       ^(\\\\_\\\\_\\\\_\\\\)\n                      ^^^ ^^ ^^^ ^\nEOC\n";

	var turtle = "##\n## A mysterious turtle...\n##\n$the_cow = <<EOC;\n    $thoughts                                  ___-------___\n     $thoughts                             _-~~             ~~-_\n      $thoughts                         _-~                    /~-_\n             /^\\\\__/^\\\\         /~  \\\\                   /    \\\\\n           /|  $eye|| $eye|        /      \\\\_______________/        \\\\\n          | |___||__|      /       /                \\\\          \\\\\n          |          \\\\    /      /                    \\\\          \\\\\n          |   (_______) /______/                        \\\\_________ \\\\\n          |      $tongue / /         \\\\                      /            \\\\\n           \\\\         \\\\^\\\\\\\\         \\\\                  /               \\\\     /\n             \\\\         ||           \\\\______________/      _-_       //\\\\__//\n               \\\\       ||------_-~~-_ ------------- \\\\ --/~   ~\\\\    || __/\n                 ~-----||====/~     |==================|       |/~~~~~\n                  (_(__/  ./     /                    \\\\_\\\\      \\\\.\n                         (_(___/                         \\\\_____)_)\nEOC\n";

	var tuxBig = "# Tux the Penguin (large version)\n#  seen when connected to irc.uslug.org\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts          .88888888:.\n         $thoughts        88888888.88888.\n               .8888888888888888.\n               888888888888888888\n               88' _`88'_  `88888\n               88 88 88 88  88888\n               88_88_::_88_:88888\n               88:::,::,:::::8888\n               88`:::::::::'`8888\n              .88  `::::'    8:88.\n             8888            `8:888.\n           .8888'             `888888.\n          .8888:..  .::.  ...:'8888888:.\n         .8888.'     :'     `'::`88:88888\n        .8888        '         `.888:8888.\n       888:8         .           888:88888\n     .888:88        .:           888:88888:   \n     8888888.       ::           88:888888\n     `.::.888.      ::          .88888888\n    .::::::.888.    ::         :::`8888'.:.\n   ::::::::::.888   '         .::::::::::::\n   ::::::::::::.8    '      .:8::::::::::::.\n  .::::::::::::::.        .:888:::::::::::::\n  :::::::::::::::88:.__..:88888:::::::::::'\n   `'.:::::::::::88888888888.88:::::::::'\n         `':::_:' -- '' -'-' `':_::::'`\nEOC\n";

	var tux = "##\n## TuX\n## (c) pborys@p-soft.silesia.linux.org.pl \n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n        .--.\n       |$eye_$eye |\n       |:_/ |\n      //   \\\\ \\\\\n     (|     | )\n    /'\\\\_   _/`\\\\\n    \\\\___)=(___/\n\nEOC\n";

	var tweetyBird = "# Tweety bird\n#  from http://pastebin.com/isRcSy01\n$the_cow = <<EOC;\n    $thoughts\n     $thoughts\n      $thoughts\n                    ___\n                _.-'   ```'--.._    \n              .'                `-._ \n             /                      `.     \n            /                         `.  \n           /                            `.  \n          :       (                       \\\\   \n          |    (   \\\\_                  )   `.  \n          |     \\\\__/ '.               /  )  ;  \n          |   (___:    \\\\            _/__/   ;  \n          :       | _  ;          .'   |__) :  \n           :      |` \\\\ |         /     /   /  \n            \\\\     |_  ;|        /`\\\\   /   / \n             \\\\    ; ) :|       ;_  ; /   /  \n              \\\\_  .-''-.       | ) :/   /  \n             .-         `      .--.'   /  \n            :         _.----._     `  < \n            :       -'........'-       `.\n             `.        `''''`           ;\n               `'-.__                  ,'\n                     ``--.   :'-------'\n                         :   :\n                        .'   '.\nEOC\n";

	var USA = "# USA flag\n#\n# from http://chris.com/ascii/index.php?art=objects/flags\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n\n  |* * * * * * * * * * OOOOOOOOOOOOOOOOOOOOOOOOO|\n  | * * * * * * * * *  :::::::::::::::::::::::::|\n  |* * * * * * * * * * OOOOOOOOOOOOOOOOOOOOOOOOO|\n  | * * * * * * * * *  :::::::::::::::::::::::::|\n  |* * * * * * * * * * OOOOOOOOOOOOOOOOOOOOOOOOO|\n  | * * * * * * * * *  :::::::::::::::::::::::::|\n  |* * * * * * * * * * OOOOOOOOOOOOOOOOOOOOOOOOO|\n  |:::::::::::::::::::::::::::::::::::::::::::::|\n  |OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO|\n  |:::::::::::::::::::::::::::::::::::::::::::::|\n  |OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO|\n  |:::::::::::::::::::::::::::::::::::::::::::::|\n  |OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO|\n\nEOC\n";

	var vader = "##\n## Cowth Vader, from geordan@csua.berkeley.edu\n##\n$the_cow = <<EOC;\n        $thoughts    ,-^-.\n         $thoughts   !$eyeY$eye!\n          $thoughts /./=\\\\.\\\\______\n               ##      $tongue)\\\\/\\\\\n                ||-----w||\n                ||      ||\n\n               Cowth Vader\nEOC\n";

	var vaderKoala = "##\n## Another canonical koala?\n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts        .\n     .---.  //\n    Y|$eye $eye|Y// \n   /_(i=i)K/ \n   ~()~*~()~  \n    (_)-(_)   \n\n     Darth \n     Vader    \n     koala        \nEOC\n";

	var weepingAngel = "#\n# Weeping Angel\n#\n# Don't blink!\n#\n# based on design found at http://shirt.woot.com/derby/entry/73182/dont-blink\n#   and http://infinitywave.deviantart.com/art/Don-t-Blink-tee-391963389\n#\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n\n                                     ...I..\n                            :XX:X\\$ . .7N..            ..\\$\\$.. .:~..\n                            X:XX.. 8XXI..         ....XX..7..7KKK8.. .\n                              N. .XXX,          ..:ZD- ..M.\\$K:XN?XX.XN .\n                              .. XX\\$.            *. .KN7XXX+ -XX,CN.,XX\n                        IXX?..                  ...--+..IXX:X:X..-ZN?DX,.\n                      .\\$XXXXX. .X               ..XX~-7=\\$7+IX5\\$...+IM+XXXX.\n                    .  +....7D=               .7=IX,: 7+..   . ,.  =+XXID..\n                    .-\\$-.. .    .             .MM-,,..... . ,=7OI.. .,:N7%\n          ..17KN. ..XX:XZ.  ..., .            .:.. .     .-IN78IN7=-.,..CMO.                     ..-.\n        ..XXXXMO..  ..8X:X= 8D0..               .8 .    .+I:N-X:XXDXXXX.ID..                  ..-XX:XX- .\n        .X:DX:XX.     ...8KI78M,                .X......-IM.D.XDXXXXXXX?:. .                .:++7CXXXXXX?.\n      . X7XXXXD8 .      .  ?8XXX.....            X+ : ...XX.X.DCO.  .+X-8X,.              . .=?2XXXXXXXXXX.\n       .=XXX887+.       .=:+?,:I.XXXXZ.         .7NO.=8+ ...M?.. 8X.\\$\\$X 7=..              ..=?X:X7++\\$NXXXXX:\n       =XX:D-...  .:..  .  ...  ...?XXX         .*. X(O)X  :X .X(O)X+X.,.                .*XXX=    . -IXXXD,\n     ,-X:XI..     ,.     ...           .       .X..........XX-....=.?N?.O?               .+XXX..       .-XX:X.\n    .:ZX8-      . .   ..+\\$:,..                .DX..D:\\$78. .XX?:XXXXD?XZ...               +DXX,    ..  .. =X:XX\n    --\\$D   .    ......+XX. .I.                  X-.*XX ..XZXXXXX.XXXXXX 8.              -.XXO.    .:=-7=..:8XX,\n  ..:7Z. .      ..  +X:XX.  X..                 :I.........XXX....XX:XX  .            ..7X7Z..  .-II\\$%>I?+,-X:X.\n   :,7..       , ...8X:XX.  XX                  .X ...  . .. .8XX:,-ZXX.              .7IXXX .  .+ODC8:\\$II7:ZXX:\n  .,=.      ..?.  ,8X\\$XX=   .XX..               ..+..  .-..%77.XXX.\\$XX .            . *IXXXX    :78DX8D08DN7-XXX.\n  .:..       *-...-XXXXX.   .,XX..                8,...O..VVVVV XX:XX               . ZOX:XI  . =ID0X0XODOX:+\\$XX..\n  .:        .? . ,IX:X:X,    .:X=.                .\\$ .N VV ....VVM+XX.              .=ZOXXX.  . IDO0X0X0D0X:0?XX?.\n  ,         ...  :IXXXXXX- .....O.                 7- O...,I.Z.X.X:X,               .7\\$:\\$\\$..   .OCO0X\\$%DCOXXD\\$7X>.\n  -  ...    ?..  =+X:X:XXDX:XX+ .. . .            ., .I. .VVVVV  XX\\$..              :.\\$XXX     +DCXXX\\$+%\\$D8ZXXDXO+\n  , .*...   .%7...,,,%%0%XXXXXXXXXXO-D\\$7J0\\$: ..-:.. ...W,XXXXX,-IXX                .:7XXXX    .CD88DCW..\\$DX:X:XX?%\n .,.-%,-.   . ?7.  .. ..+%8XX:X:XX:X:X:NO\\$DX..,XXX.. . ..XXDXXI\\$XX .               .-O:XX     -CDOI8XX*.*DX:X:X7=%\n.--.:\\$.+ .     .I%       .:7XXXXXXXXXX    . . ....      ..*.-XDXX..               .-:XXX .   .=78CD8XX+  +IXXX>..Z\n ?=:=?.+.        .XX..      ...?8:XXXX-.           .    .:-Z+XX- ....... .       .+,DXX.    .7.,\\$88DXX:  .:OD%..%O.\n ?+IO--+           -XX.  .        .-7O.    .              .\\$\\$.. .7.M.\\$XXX:...   .++DZM.   ..I\\$  .\\$8DX7.,..  .. IXX..\n +??D,\\$\\$              .X8.       .        .*      .   . ..  ...*XX X,X:X-: ON?...,XIO     .-O\\$.   .-. =O,. .  OXXX-.\n.:?ID,Z7              .. ....:I-..        .      ...  .+..   :7 XXD. XX,:...+XXX-.,*..    .=\\$?..  . .:O8.\\$D7 .IXXX?\n ,+78:Z7 . ,*. ..*.           .....  .....D . ?+.     .:+.  -X+.XX7.8XZ.*.  ..+7XX,     ..*+:.  ... ,8:+.-D0..7XXX7.\n.,=\\$8I88   ,7  :XD-.  ....             +.:,...--|.. ..      7% CXZ..XX,.*.  ....+XX+.   ..:.    .\\$O..DO,  8D+.-XXXX,\n :-\\$X0XX ..,*. -DX,... 7... .           .X...::::::-O...    .  .. .7XX...  .,XX. 78-+:..        .OZ.-87   CD\\$.7XXXX.\n.--7O0XI.. :-O:ZM8...-*O.  -.           :% .::-,,::..:,::::--*+I88XXXD....X:X:XX.7XXXXX:        .8%7DO....-8O77DXXX,\n -=78XX.   .**%XX .   %7:..Z .          8. .:X7.::*,* :,,,::-:-*7\\$DXX+.,.XXX?I+.7X:X:X.-X:      .-%%8O .  .8ODIDXXX,\n -=7IX8.    .*87.   ..\\$O7-I\\$..    -.    .... .. ....:,:O:,::-: .7IDXX ..XX..  .:X:... XXXXX.    ..=7+.    .=0D7CXXX:\n.-+7IN..          . :..O\\$7X,      =.  ..,.  ..   ...   O..7.::. ?.XXX,.  .      ..  .X:XZ08. . .          ..\\$D8%XXX:\n.-7II.            .+%. :77:.    .:..  . XD  .-  :...  .:. .   ..... ..          .-  ,X7...,X%.:            . DO8XXX,\n ,I* ,       -?:..8XX. . .....   ...   :XX...X .\\$  , .?%: ,? ==7X7., .           D..? . .XXXXX:%...         ..%8DXX\n ..      ..  DOX .8XX..-..:==.:  . .   .X%  .X .% .%. XXO DX.D?\\$X?               XX.  . 7XI.XX.. :.   ....    ..-?,.\n        ... .DXX..8XX .?8-Z\\$?..*.* .    X-  .X..7 7\\$..XXX.DX+XXXX,   .        . .XX.  .  .  .XX...     .*.  ..  ..\n  ....  :,:. IXX-.8XX .-XDXX- .++\\$..  . X,  .X.%,.+% .XXX.\\$XIDDDD. .=.  .. ...,8XO..        -XXXX ..    *...., .....\n   ,,\\$. -:*. -DXO+8X\\$ .,ZXD7,  .:...  .,X.  .?.X- +-..XXX 7XOXXXXO. ..  .NXXXXXXXXX         87%:XX..    -..I7\\$:.,-.\n . .DX,.-I% ..\\$XX8DXI   *OI-..  ..    .XX   ..7X8.?,  XXX.\\$DX:XXID     .XXXXXXX8.. ..       .+XXXXX+  ..,  \\$XX..:?..\n . .8X-.-%8...*DO\\$XX7   ......         XX   . XX .I.. XXD 7XX:XX?,     ,..              .8D\\$X:XXXXXX   ... \\$8X. ,I7.\n    I8\\$.-OXX..,\\$D8XX... .  -,,....     XX.  ..X  .?:. ZXO -XX:X\\$,:...                 .... :XX:IXXXX   ... %8X  .\\$X.\n..  -8D.-O8X...-\\$CX+ ..D,  -X7,-,.     XX   . 7. .I7 .7X\\$ :X\\$XX-.XXI .                .:: ..X:XXXXXX.. ,. .%XX  .CX:\n..  .CD.+OCX  . .,. .-.XO. ?8I,IX.  . :X7.  ...   78..8X-.+X+XX, XXD.     .     .:    ..*N= -XXXXXXD  .:. .OXN  .8XI\n .   78*COyX,       .8*XX  =8....   ..\\$X..... .   XXI XX .XX.XX..?XX..7.  ?,    ..      ..?XXXXXXXX.  ,:. ,OX*  ,DX\\$\n .  .%D7CO7X\\$ . .... X7XX   ..        XX . ..   ..XX\\$ X8  XX.XX...XXX%8   %-              .. .=\\$XX..  ,:. :ZX   :DX%\n.... +XO8O7XD ...+.  X8XX.          .:X.  .7-    ,XX\\$.D= .XX 8X   DXX\\$X:  XD.                         .-,.*OX.  =DX%\n .,. .XD8O\\$XX  ,,\\$:. OXXX-          .XX.  .D?.  .XXX:-8. +XX DX:   XX7X\\$  XX                          .=-IZXD   +DX8 .\n ::.  XXXXZX\\$  .,\\$D  7XDXI          OX... O8?   .XXX.%8  IXX ZXI   ZX7XD. DX.             ..   .      .=+Z\\$X=   7DX%\n.-:.  IXXN7X. ..:IX. 7DNK\\$        ..XZ  .\\$XX- ..OXX8 %Z  7XX.-XX ..+XDXX+.7X-.           :+   .D      .-*\\$+X... %DX%\n.*,,. .XXDIX, ..??X. IDXXI        .XX   .XX8,  :XXX  %% .=DX..XX.. .XX:X8 .XX.           .=.  .X      ..+??X  ..%DX?\n.+*-  .OX7?%. .:%77\\$.7DXX         -X7.  XXX+...XXX:. ZO...DX\\$.OXX  .+X\\$XX,.XX,          ....  ,X       .-8\\$.  .-OXX+\n *I*,. .8+ZI  .=D+XX ,8XX.        XX. .8XX%...7XX7...%X   8XX .XXX  .X\\$XX\\$.XXX..              -X      . -X. ..?+DXX:.\n =O7%:. :-I  ..*O+XX  I8,.       .XI. .XX%.. -XXD   .:X:  .XX* ZXXZ..\\$X7XD.:XX7.              ?8        .-...I8\\$XXX..\n -X\\$O%:    .,..=7-XX. .,..      .*X   .OD.  .:+%..   *XD  .8X- .7XX..XX..II +CO               -.     ..   .,%XODXXX...\n :XOOX-.  ..?..:I:XX..                 .\\$,  .:.     ..,:  .I?,. +8?..\\$%.  . ..,.            . ,.    .?..   =XX:XXXX\n .XXXD=    ,O. .I.\\$M=                    .     .     ...  ......-....,..     ..             ..X=     X .   -DXX8XXX\n .XXX8*.  ,,\\$:. .?+XX                                               .                        =N?..   D     :OXD8XXX\n  %XXXI.  ,,7\\$. +I78N.            ...                                                        ID?.   .I    ..%X8CXX+\n  +ICOI   .,:Z- .-XXD.            .,..                                    ..                 7D*..  ..      \\$D\\$ZDD..\n  ,IOO+.  ,::Z7..7\\$\\$D              ,.. .                                  .@.                -O\\$.:  7.     .78OZ8D..\n  .888..  ,-:7O+.7778.            .,.  I .     ..      ....     ...:..  . OO..               .--*:. 7.. .   :%XZ8.\n.,.:%I..  ,+:?O?.=I7\\$.             .   .      ....... . ::    , .*7%..  :.OO=.               ..-+..,?.. .   .I?OI.\n    .     ,7:=\\$7,,II7.            ..  ..   .. .-    ..  +7.   :  77\\$.-..? %%+                 .7\\$. +=:,..   .,..\n          .+::7Z7.:I+.            ..  :.    .. -* .. \\$  I\\$ .  -,.I\\$7 =  7.XXI                 .*...*,,.... .    .\n        . .-*:I\\$?..*.              .  *.    .  -7.. .7. I\\$.,  :-.7\\$7.=  7-\\$\\$7                ..    .,- .  ..  ...\n    . ...  :7:777. .                  -.    +- +7.   7..7\\$+ ..?.?77  =  7+777                 .. ..7?I....-   .. .\n    = .....,7:-7I-.                   -.    II =I:   :?.?I7   ?.?II. = .I?7::                   ..+--?.. -*,  .=?\n    ,,..,. ,**,I?I..              .   ,.   .II :*?.   ?.??I.. ?.:?I. - .I?II?                   . ,:*?-. ==I..=I?.\n    .=...+ .*+,-??.*              .   .:.   +7:.-7.  .+,+??:..*:.+7. :. 7777+                   .,7+++7. -+,.-=7*.\n    .=...*. :=:,++:.  .           .   .:. .. =+=.-*: .=*=**= .:-.**, ,..*****.                  .,:=*=*.,:*..***=.\n    .....-. ,*-,===.              .   .:  .  ,== :-*. .,*,*+*. .-.** :  -****.                   ,:**:*.,:+..*+*:\n    .. .... .:-:=== .             .   .:     .*: ,:*.  .*,--*,  -.-***  :****,                   .:---: ,:-..-**,\n.   ..    .. ,--:-- .             .    ,     .-: .,-,.  -::--,..:,---:, .----:                   .:=:=. .::.,===.\n       :-::. .::::- .             ..   .     .::  ,:-   .::--:, ,,,;:.  .---::                   .,:.: ..:,.::-:.\n     ,..,,,.  ,::::                .   .     .,,  ,,:.  .:,:::,..,.:::. .:::::                   .,,,, ..,.,::::\n        ,,,,  .,,,,                .   ..     .,  .,:,   ,,:,,,. ,.,,,.  .:,,,                    ..,.. ...,,:,,\n    ,,...,,,   ,,,.                    .      .,. .,,,   .,,,,,. ,.,,,.. .,,,,                    .    .,,.,,,,.\n     .,......  ...                     .      ...  ,,,   .,,,,,. ...,,,. .,,,,                          ...,.,,.\n      ,......  ...                  .  ...     ..  .....  .................,.,.                           ...,.\n      .......   .                        .     ... ....   ....... ....... .....                      ...........\n      ... ....                          ..     ...  ...   ....... .............                        ........\n       .......                          ...     ..  ...    ...... .............                       ........\n        ......                          ...     ...  ...   ... ..  ............                         ......\n        .....                           ....      .   .      ....  ............                        .... .\n        .. ...                       .. . ..           ..       .   . .... .. ..                       ....\n           . .                          . .       ..   ...   ...      ... ....                          .  .\n           .                              .       ..          .       .  .   ..\n             .                           .                    . .    ..    . .                             .\nEOC\n";

	var whale = "# whale\n#\n# modified from https://www.reddit.com/r/pics/comments/25ji0n/man_face_to_face_with_whale/chi1kdy?context=3\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n     $thoughts\n                '-.\n      .---._     \\\\ \\.--'\n    /       `-..__)  ,-'\n   |    0           /\n    \\--.__,   .__.,`\n     `-.___'._\\\\_.'\n\nEOC\n";

	var wizard = "# Wizard\n#\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n                     _____\n                   .\\'* *.\\'\n               ___/_*_(_\n              / _______ \\\\\n             _\\\\_)/___\\\\(_/_\n            / _((\\\\- -/))_ \\\\\n            \\\\ \\\\())(-)(()/ /\n             ' \\\\(((()))/ \\'\n            / \\' \\\\)).))\\\\ \\' \\\\\n           / _ \\\\ - | - /_  \\\\\n          (   ( .;\\'\\'\\';. .\\'  )\n          _\\\\\\\"__ /    )\\\\ __\\\"/_\n            \\\\/  \\\\   \\' /  \\\\/\n             .\\'  \\'...\\' \\'  )\n              / /  |   \\\\  \\\\\n             / .   .    .  \\\\\n            /   .      .    \\\\\n           /   /   |    \\\\    \\\\\n         .\\'   /    b     \\'.   \\'.\n     _.-\\'    /     Bb      \\'-.  \\'-_\n _.-\\'       |      BBb        \\'-.  \\'-.\n(________mrf\\____.dBBBb._________)____)\nEOC\n\n";

	var wood = "#\n#  木\n# 木木\n#\n\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts   --木--\n       ／｜＼\n     ／  ｜  ＼\n  --木-- ｜ --木--\n  ／｜＼    ／｜＼\n／  ｜　＼／  ｜  ＼\n    ｜        ｜\nEOC\n";

	var world = "# World\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n          _,--',   _._.--._____\n   .--.--';_'-.', \";_      _.,-'\n  .'--'.  _.'    {`'-;_ .-.>.'\n        '-:_      )  / `' '=.\n          ) >     {_/,     /~)\n  snd     |/               `^ .'\nEOC\n";

	var www = "##\n## A cow wadvertising the World Wide Web, from lim@csua.berkeley.edu\n##\n$the_cow = <<EOC;\n        $thoughts   ^__^\n         $thoughts  ($eyes)\\\\_______\n            (__)\\\\       )\\\\/\\\\\n             $tongue ||--WWW |\n                ||     ||\nEOC\n";

	var yasuna_01 = "##\n## やすなちゃん\n##  \n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n  \n            . .: ───:. .\n         .／.: .: .: .: .: ヽ\n        .:   .:l.:   .: .: .:.\n        |.l:..ﾊ.ハ..|ヽ.ﾄ､:: |\n        |:l.:/ヽ､_ヽ|_ノV:.:.|\n        |:lﾊ.  {j    {j  |:ヽl\n        ﾉ:l} ''        ''|:ノヽ／ )\n        ヽﾍ:ヽ.､ r---､  ｨﾉ ┬' `／\n     γ::ヽ  ｀^Y`TﾇΤ` {__├'`'\n     ｀‐< ＼_ ハ |:|  Y\n          ヽ_>､|  |:|／|\n               /   V   l\n             〈        〉\n           〈:｀-:';`-´:〉\n            .>-:ｧ─--‐r-:ｨ\n            /  /     |  |\n           /  /      |  |\n          /-,/       |--|\n         に7         |二|\nEOC\n";

	var yasuna_02 = "##\n## やすなちゃん\n##  \n##\n          $the_cow = <<EOC;\n           $thoughts\n            $thoughts\n                           _.. .:-―-:. .._\n                      .: .: .: .: .: .: .: .: .: \n                   ／ .: .: .: .: .: .: .: .: .: .＼\n                 ,'         ,!    ∧           : .: ヽ\n                /, .:: :｜./ |.:./ヽ.:iﾍ.: .: .: .: ::.\n               ,''|.:: .人/--|':/  ヽ:| ､＿.: : .: .::|\n                  |.:: ｲ  ,,=､ﾚ        ゞ=ﾐ､.:|..: .: :|\n                  |.:: ｜{{    }}     {{    }}八.: .: :|\n                 /.: : /  ゛= \"        ゛= \"    ;.:r ､:|\n                /,.ｲ.:〈                     ,, //' ｝:|\n               '  ヽ:: ゝ、        ｰ--┐       //  ノ::.\n                    ヾ::.､＞ .    ヽ _ﾉ    ..  ＜¨ｨ.:}~＼\n                      `゜ヾ/｀>了、.    v 〔:／|:/  レ'\n                        _ . -/: ,K:::>､/: :ﾄ._\n                       |: :く_.:|/:〈 /: :}: /~ヽ\n              r「「「ｈ,>:|: <: |'::ｿ::<¨.:n｢「「!､\n              ゝ＿_ﾉ /: : |::ヽ |::/: ／: :.ﾍ＿_ノ｝\n               | ￣ |,': :/: : ヽ:' ／ : :.:| ￣ |:}\nEOC\n";

	var yasuna_03a = "##\n## ソーニャちゃんおめでとー！ソーニャちゃんおめでとー！\n##\n$the_cow = <<EOC;\n $thoughts\n  $thoughts\n   $thoughts\n            . .: -----  .\n         ／: .: .: .:.: .:＼\n        /    ..  . l.: .: .:ヽ\n       : .: ,/|-/|:ハ.:|-.ｌ.:\n       |: :ノ |/.|/  ヽ|.Vﾊ.:|\n       |.::|  =＝     ＝= }.:| \n       |.γ|| ''  ＿_   ''{::ﾊ \n       ﾉノﾊﾘ   ｛   }     ﾉV\n       ∨Vvヽ､._  --'_ .イV\n            γ:/:{.又 }ﾍヽ \n          ／:〉:V ﾊ.ﾘ〈: ＼ \n        ／ : Vヽ:V// /:V ::＼ \n    rイ: : ／|: :＼Vノ: :|ヽ: ヽ-､\n   ｢  ヽ:／  |: o :  o:|   ＼:/ 」\n    ー'    ./: : : : : ﾊ      ー' \n          ./::o: : : :o ﾊ \n          /ヽ: : :Λ: : :ﾉ:、 \n        〈:::￣￣:::￣:::::〉 \n          ＼:__:::::::__:／ \n            |  Τ￣Τ | \n            |  |   |  | \n            |''|   |''| \nEOC\n";

	var yasuna_03 = "##\n## ソーニャちゃんおめでとー！ソーニャちゃんおめでとー！\n##\n$the_cow = <<EOC;\n            . .: -----  .\n         ／: .: .: .:.: .:＼\n        /    ..  . l.: .: .:ヽ\n       : .: ,/|-/|:ハ.:|-.ｌ.:\n       |: :ノ |/.|/  ヽ|.Vﾊ.:|\n       |.::|  =＝     ＝= }.:| \n       |.γ|| ''  ＿_   ''{::ﾊ \n       ﾉノﾊﾘ   ｛   }     ﾉV\n       ∨Vvヽ､._  --'_ .イV\n            γ:/:{.又 }ﾍヽ \n          ／:〉:V ﾊ.ﾘ〈: ＼ \n        ／ : Vヽ:V// /:V ::＼ \n    rイ: : ／|: :＼Vノ: :|ヽ: ヽ-､\n   ｢  ヽ:／  |: o :  o:|   ＼:/ 」\n    ー'    ./: : : : : ﾊ      ー' \n          ./::o: : : :o ﾊ \n          /ヽ: : :Λ: : :ﾉ:、 \n        〈:::￣￣:::￣:::::〉 \n          ＼:__:::::::__:／ \n            |  Τ￣Τ | \n            |  |   |  | \n            |''|   |''| \nEOC\n";

	var yasuna_04 = "##\n## 湿布\n##\n$the_cow = <<EOC;\n            $thoughts\n             $thoughts\n              $thoughts\n                    .  .:----.:  .      \n                  ／    .: .: .: .:＼\n                 .          .. .: .: ヽ\n                /.: :/|:/.:ﾊ::ﾊ : .: .:.  \n              ノ.: ./-|/ |/ V- V､.: .:.:|            \n               ｜:ノ   _     _   V: .:,:|\n                |:}  =＝     ＝= |:l､:.:|\n                |:ﾉ''    ＿_   ''|:| }::|\n               八:ヽ.   V_ 丿   .|ﾉｲ: :八\n                 ヽ/≧=-z:-:r:=≦l:ﾉ|:／\n                    ／/ ﾚﾇﾘ: 〉 ＼\n                   / :〉|/l/:< : ハ  \n                  /:}:{:|:/:/ : : :.     \n                 /: { : ': /: : {: :|            \n                 {: : :ヽ:/: : : : :}          \n                /} :}:o : : o: {: : ﾊ                \n                {: :ﾘ: : : : : |: : }  \nEOC\n";

	var yasuna_05 = "##\n## ダーツの才能\n##\n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n       ／ .: .: .: .: .: .: .:  .: .: . ＼\n     ./   .: .: .: .: .: .: .:  .: .: .: .ヽ\n     /          /  . ..l..  ヽ.: .: .: .: .:.\n    ,    .. .: /  .| : ハ: .|  ＼.: .: .: .: .\n    |.: .:.l.:/  ヽ|.:/  ､ .|.ノ ＼ .l:.: .: |\n    |.: .:.|:/.ｨ≠ﾐ|:/    ＼| ィ≠ミ､|.:.: .:|\n    |.: .: ノ /Y::::ヽ       Y::::ヽヽ＼ .: ｜\n   /:.: /^|:|{.{:::::}       {:::::}.} |:|ヽ:､\n ノ:ノ: { |:| Ｕうーソ       うーソ  |:| }ヽ:＼\n    | : ヽ|.|  '' ￣           ￣ ''U|:| /:|\n     :: ::人|                        |人::ﾘ\n     Vﾊ:: :: \\                     /::  ﾊ/\n      \\|ヽ:: ::ヽ､     --      ,イ::／|／\n          ＼| ヽ:≧=r-r---r-r=≦:ノ|／\n             . :´:.ヽ二二.ノ: :｀: .\n            ／: : :  ／ハ＼: : : : ＼\nEOC\n";

	var yasuna_06 = "##\n## 策士\n##  \n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n \n              ,.:￣￣￣￣:.､\n            ／ .: .: .: .: .:＼\n          ／   __ l  __    .: .ヽ\n        ／:  ./ＶＶＶ＼:ﾄ､`.: .:.\n   (＼  ￣/.:ｲ.ｨ=ﾐ` ´ｨ=ﾐ､ヽ,: .:|\n   {ﾐと^ヘl.:ﾉ{ぅｿ,  ぅｿ}|:|^ヽ:|\n    ヽ〃: ﾚ｛    __ -､ \" |:| ﾉ::|\n      ＼: :ﾊ＼  {    }  ,ﾚイ:::八\n        ＼ :V.:>:ニニﾆ:＜ＶＶＶ\n          ＼ :v:〈|父 /:|:ﾊ┐\n            ヾ{:｢|/:|/ <:/:ﾉ＼\n              {:\\|::/:／:{: :ヽ\n              {: : : :`: /> : :>\n             / : ﾟ : : ﾟ : Y: :/\n            /: : : : : : : |ヽ/\n           〈: :ﾟ: ﾊ : :ﾟ: ∟ｺ\n           /::---':::------く \nEOC\n";

	var yasuna_07 = "##\n## ごぼう\n##  \n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n     $thoughts\n               ＿＿＿＿＿\n           .:´.: .: . : :. `  、\n     ..: ／.: .: .: . : .: .:   ＼\n    .::／:::       ﾉ   /､         ＼\n   ..:/.: ::.:|＿／::|:/  ＼:__|:  .\\\n .:: :::: :::/|／｀ヽ|/    '＼:ﾄ、:  .\n .:::|.:: ::/:ｨf于ミ     .ィ≠ﾐ､Ｖ: :. .\n..:::|.:::ノ::{{:::}       {:::}}{: |＼|\n..:::::::_::|::うﾆソ       う:ソＶ: |\n.::: /.:/ |:|:ヽヽ       ｀      }: |\n.:::/ｲ:{  |:|:    ／￣￣ ｧ      ﾉ  :|\n ..::|.ゝ,ヽ|:   /      /     ／:::八\n .:::Ｖ:::::＞:._ヽ、 ./__ .イ:ﾊ:／\n  ..::＼|＼:斗:ｰrﾍ`ア又＜Ｖ|／\n   ..::::／⌒: :|:ＶＶ{ヽ:＼\n      .:/.: :|::l::ﾍ}/\\|:}:.＼\n    ..::｢.: :|::＞:Ｖ//|〈:.}.}\n  ...::/.:: :|::＼: Ｖ/| / :}:.┐\n ...::/.::::rｰ::::＼:Ｖ|/〈::::.ヽ\n..:::/.::::ｲ::::::: ＼ Y::ヽ:::::.＼ \nEOC\n";

	var yasuna_08 = "#\n# ごぼう2\n#\n\n$the_cow << EOC;\n  $thoughts\n   $thoughts\n    $thoughts\n          ,.:──‐-:.,\n        ／:.           ＼\n      ／:. :. :. }:. :. :.ヽ\n     .: :. :. }.:/＼.:|,:. :ﾍ\n     |:.:. :. /Ｖノ ヽﾄ＼:|､ﾍ\n     |:.:. /Ｖ_ﾆ    ﾆ＿_ {::ﾍ\n     |:.ﾍ .| ΓT      | |Ｖ.＼\n     |:{ |:|.l｜      | |八:ー\n     ハ:`:Ｖ､l｜∠二l.|.ｲ:ﾊ:ﾉ\n      _Ｖ＼;＞=r rr r=＜ﾊ／   ＿＿\n     |ざ |ﾍ :{Ｖ/V:}:＼      |ご  |\n    {ﾐ}く{)}:＞Ｖ/< :  ＞-:-'{}ぼ{ﾐ}\n     |ろ_|:ﾉ:＼:Y / }ﾐ : : : |  うY\n          ﾉ :o: : :oj `ー─-´ ￣￣\n         / : : : : :{\n        /: : o : :o:ﾍ\n      〈 : : : /\\: : 〉\n       /::ー──'::ー‐ﾍ\n     〈:::::::::::::::::〉\n       ￣|￣｢￣|￣|￣\n         |  |  |  | \nEOC\n";

	var yasuna_09 = "#\n# ちびきゃら\n#\n\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n    $thoughts\n    \n           ____\n       ,: .: .: :.ヽ\n     ,'       /\\   ｉ\n     {: .:ﾉﾚﾍ/  Viﾍ:}\n    .{,､〈 Ｏ   Ｏ{.:.\n    ノヽ\\!\"       }.:ﾊ\n      Ｗﾊw=-､へ,ｬ<,V'      \n         /ﾍ }{./\\\n        ;: i:V:!;}\n        |:｜: :｜}\n        |:|:｡: ｡l}\n        >-'-ﾟ-'`ﾟu\n        ｰi-i～i-i~\n         |.|  |.|\n         |-|  |-|\n         ヒｺ  ヒｺ \nEOC\n";

	var yasuna_10 = "#\n# 何でそういうときだけ凄そうなの！\n#\n\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n    $thoughts\n             ＿＿＿＿\n       ＜ :: :: :: :: `丶､\n       ／   _, ｨ:ﾊ ､＿: ::＼\n     ∠:: :/ |/|/ \\/  \\/:: |\nrヘn  /:\\/ c=＝.::.＝=っ\\/ |  rvへ\nヽ／＼i:｜   ┌──┐   i::|／＼ノ\n  ＼::|(||   |:::::|    ||)|::／\n    ＼|人|.、|:::::| .ｨ|ﾉ:八／\n      ＼\\/\\/>|:::::|<\\/\\/／\n        ＼ :::>TﾇT<::: ／\n          Y : ＼W／ : Y \nEOC\n";

	var yasuna_11 = "#\n# きゅーっ！\n#\n\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts          . .: -ーー― :._\n    $thoughts       ／.: .: .: .:     ＞  r⌒ヽ\n           / .:         ｜.､.:＼  ﾉ ノ\n          .: .: .:|＼  |斗ﾍﾄ.:.:Ｖ  /\n          |: .: /\\|ノ＼| ／ Ｖ::Ｎ./\n          |: .:/ c─-        Ｙ:| /\n          |:ﾊ:{``   ,  --┐  人V /\n          ﾉ:L＼>   く_,￣┘／  ＼\n   /⌒￣￣￣|￣￣＞--r-rｭ＜|   ／\n   L_,vー─-|    ､ }  ＶYﾊ   Y\n             ￣￣Ｖ  ｜/∧   ﾍ\n                  {   |//∧  ﾍ\n                  {    ＼//   ﾍ\n                  {            ＼\n                  ｝             >\nEOC\n";

	var yasuna_12 = "#\n# からあげ\n#\n\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n         $thoughts        .:  ￣￣￣￣:.丶､\n               ／.: .: .: .: .: .: ＼\n              /    ／|    /\\.:| .: :.\n             / .:|乂 |/{:/ _乂/\\ .:.:|\n           ノ.:\\/ｨ庁ﾐ` \\/ｨ庁ﾐx  \\/:.:|\n             |:}{弋.ﾉ    弋ノ } /.:.:|\n             ﾚ:ﾘ''          '' ｜:ハ:＼\n             {人       ,、    ,｜/ノ:厂 \nEOC\n";

	var yasuna_13 = "#\n# 転んでも泣かない！\n#\n\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts    ｡\n         $thoughts       ＿＿＿__{_    o\n      ○   ヽ￣    .: .: .: ｀丶＿_\n           ／ .: .／;|.:/|::∧.:＼        О\n     (ヽn∠ .::|∠二:|:/:|ン-:∨:「ﾚ^L\n     ζ, ヘ /::(___)|/:(￣￣ )|.:/、  ζ\n     `く: :/_:ﾉ(_) ＿＿ ￣(_) |:/:`ｰ:/\n  。    ＼|＼人   |/   `⌒ヽ ｜/ : :/  o\n      __┌   ＼`ヘ|/ヽ/ヽ／^ヽ/／/:/／/\n      ＼                        /:／ / \nEOC\n";

	var yasuna_14 = "#\n# くっ、くぅー！\n#\n$the_cow = <<EOC;\n       $thoughts\n        $thoughts\n         $thoughts\n                 .:-────-:.  .\n             .: .: .: .: .: .: .: :.\n          ／.: .: .: .: .: .: .: .:  ＼\n         .: .:          /:  |        :.\n        .: .: : :  |.:./ |.: ハ.:.|::|.: :＼\n       /.: .: .: .:|.:/ u|.:/ u ､:|::|.:｜―`\n      /.: .: .:|:,|.＼._ｨ.:/ ､_／｜∨,::|\n     /.: .: .: |:/ィ≠ミ |/  ィ=ミ､ ∨::｜\n    /..: ,--|:｜  {んi:i}     ri:i}} ﾊ::|\n   /.ノ.:/へ|:|.   ∨:タ...::.ヾ:タ  .:.:､\n   ／:: :ﾊ (|:| u ''       '    ''  {:|ヽ:＼\n     {: :＼_|:| ｕ   __          u ﾉ:｜\n     ∨ﾊ:ﾊ:ヽ.|､   （- `ｰｧ      ..ｲ::/ﾉ\n       ,:＜:￣/|､＞:._￣..:-=≦::ﾊ:／\n      /: ヽ::/:| ＼_ィ .ハ＞:、\n     」: : :く:｜／{;;}∨: }::ﾊ\n    /:＼ : }/￣`Yヽ:∥:／: /:「Y二ヽ\n   / : : : /  ￣}-':/::〉: }:/Y{─ }\n  /: : : :/  .二ﾌ::/::/: : ﾘ::ﾊ{-- ﾉ\n./へ──‐ﾊ  ,-ｲ :/::/ : :ﾑ:-{､_エノ\n{: : : :.ヽ>イ:|:/::ノ: :/ : {{ ／ﾉ \nEOC\n";

	var yasuna_16 = "$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n         ..: ￣￣￣￣: :.\n       ／::  /｜.:/ |.: .:＼\n      ,  /｜/  |./  |.ﾊ.: .:ヽ\n    ./.:ｲ__ノ   ヽ､___∨.: .:.\n   ./: .:≡≡     ≡≡.|.: .:｜\n   /ノ|/} }.      } } |:ﾊ:.:｜\n     .ヽ{,{ -~~~- {,{｜:/ﾉ:从\n      ∨v､＞z-r-x-:r＜/ﾚﾚへ \nEOC\n";

	var yasuna_17 = "#\n# さっそく試してみよう 道具持ってないから作るしかないかな\n#  \n\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n     $thoughts                 ____\n                 .: :<::. ::.>: :.\n               ／:: ::. :. ::. ::`:、\n               `::. ::.ィ:.i::.、::.ヽ\n             /'      ./|..ﾄ.}V.. .. ﾊ\n            '.. .. ./L/｜:| 一V::. ::１\n            i::. ::/}/` V:| V Vﾄ::. ::i\n            |::. :/Y芋ミV!Y 芋ミ|::. .|\n            ,::. ハ {::}  V {::}}:r,:代\n            /::. :}  つﾉ    つﾉ｜:レ:}ゝ  ヽ\n              V::八    r一 ┐   ｨ!::.:ﾘ      }\n       ｛r     ＼ﾊ:＞- .一-'.s<:ハ}ヽ}   __ノ ﾉ\n        弋二一   ヽ:{＞}_ノ  / ゝ､\n                ｡＜   〈ﾊ〉  {    `、\n              ／     i       `､.    `、\n            ／    フ^|   　   ',ﾞ、   `、\n           く   ／   |         ', ﾞ、y ヽ\n           tゝ_r     r          ',  ><一'\n                    /  ゞ＿      '\n                   /      一      `\nEOC\n";

	var yasuna_18 = "#\n# ま、ありがちな言い訳だよね\n#  \n\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n                      ,:二二二二:. .,\n                   ／.／＿＿＿_  ＼.:＼\n                  /. /／.: .: .:＼  : .:＼\n                 /.: .: .:/｜:/\\ .:＼}.: .:.\n                .: |.:/一/ |:/ 一.:}: .: .:｜\n                |.:|ノ |/_｜/ _  \\/ﾍ: .: .:|\n                |.: ｜= ＝    ＝＝= \\/}: .:|\n                |:: ﾘ''           '' /:/､.:|\n               ノ:|:人    一一 ､    /:/ ﾉ.:|\n                , ┴＜＼  {     ｝ ,{:/イ::八\n               /_..   ＼` ー┬一r＜:八八／\n               ／  T＼   `＜}ゞ=彡'⌒＼＼_>\n              /___ |  >､    ｀''＼   ｜\n             /ﾆ}::\\/／  ＼       ｜  ｜\n          　{ﾆﾉ:: /''＼ | `|r--ｯ＜|_／|\n           /__   V    ｝|  》=《      |\n           ＼ ＼/｀一ﾍノ|  { 6 }     ｛ \n             ￣        ｢   ゞ= '      }\n                      ﾉ               〉 \nEOC\n";

	var yasuna_19 = "#\n# やすなちゃんのまんまるお目目\n#\n\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n\n:. :.孑|:/仔:./  ＼:.| V｜:. ﾄ:. :.\n:. :/  |/  |:/     ヽ|   \\/:!\\/:.:\n:. / ,ィf芋ミ     ィf芋｀:V  .\\/.\n:./ ,' :'::::ﾊ      ,':::::ﾊ ヽ /:.\n:t  { {k)::::!     !k)::::!  },'.:\n:ﾊ    弋 一ソ      弋 一 ｿ ,: ::\n:.{      ￣    ,       ￣  ; :./\n:.| ''                  '' |:./\n:.ﾄ､      ` ､      ノ     ﾉ!:/ノ\nﾄ､!:＞ ､.     一  '   .,＜:|/::.\n:: :: :: ::>z-一-z<:: :: :: :: :.\nV|＼:/}ﾍ/  `ー又ー' \\/}ノ{／|:／\n  ,z'￣ ﾍ   /{ .ﾄ､  /￣  ヽ\n／      /\\./x 一 ﾐ./       ＼ \nEOC\n";

	var yasuna_20 = "#\n# yasuna_20.cow - もしかしたら新種かも！\n#\n\n$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n    $thoughts            ________\n             .:          :｀丶\n           /.:   ｛ :｜､  .: .:＼\n          /   |.: /\\.:|ﾉ＼.} .: :.\n         .: .:/\\乂  ＼ｨ=ミV.:}.: |\n         |.:\\/ ｨ=ﾐ    ﾋソ｝V:|.:｜\n         |.:ﾊ{ ﾋソ '    ''｜:|ヽ｜\n         |.: ﾊ''          ｜:ﾉノ:＼\n        丿.:|人    ⌒ヽ    ｲ::\\/ ￣\n    /^^ﾍ  \\/Vv:＞=rr::rr＜vV\\/\n  ｛   ﾉ    ノ   \\/ヌ\\／ ＼\n    ＼  ＼,く  }   |:|   V ＼\n      ＼     >ィ   |:|   ｝  ﾉ\n        ＼／  ﾉ    |:|   }-く ＼\n             /      V     \\  ＼  ＼ \nEOC\n";

	var ymd_udon = "##\n## 山田うどん\n##\n$the_cow = <<EOC;\n   $thoughts\n    $thoughts\n  \n             _ - ￣ - _\n           _-_＿＿＿＿_- _\n         ￣ｌ  ●   ●  l￣\n            ヽ､_ ⌒ _ノ\n         _ -‐ニ ￣ ニ‐- _\n  /⌒ ‐ﾆ‐ ￣   /    \\ ￣ ‐ﾆ‐⌒ヽ\n ヽ､_ノ       └-ｕ‐┘      ヽ､_ノ\nEOC\n";

	var zenNohMilk = "$the_cow = <<EOC;\n  $thoughts\n   $thoughts\n    $thoughts\n\n     iﾆﾆi\n    /   /ヽ\n   ｜農｜｜\n   ｜協｜｜\n   ｜牛｜｜＿\n ／｜乳｜｜／\n ￣￣￣￣￣\nEOC\n";

	function convertToCliOptions(browserOptions) {
	  const cliOptions = {
	    e: browserOptions.eyes || 'oo',
	    T: browserOptions.tongue || '  ',
	    n: browserOptions.wrap,
	    W: browserOptions.wrapLength || 40,
	    text: browserOptions.text || '',
	    _: browserOptions.text || [],
	    f: browserOptions.cow,
	  };
	  if (browserOptions.mode) {
	    // converts mode: 'b' to b: true
	    cliOptions[browserOptions.mode] = true;
	  }
	  return cliOptions;
	}

	function doIt (options, sayAloud) {
	  const cow = options.f || DEFAULT_COW;
		const face = faces(options);
		face.thoughts = sayAloud ? "\\" : "o";

		const action = sayAloud ? "say" : "think";
		return balloon[action](options.text || options._.join(" "), options.n ? null : options.W) + "\n" + replacer(cow, face);
	}

	function say$1(browserOptions) {
	  return doIt(convertToCliOptions(browserOptions), true);
	}

	function think$1(browserOptions) {
	  return doIt(convertToCliOptions(browserOptions), false);
	}

	exports.ACKBAR = ackbar;
	exports.APERTURE = aperture;
	exports.APERTURE_BLANK = apertureBlank;
	exports.ARMADILLO = armadillo;
	exports.ATAT = atat;
	exports.ATOM = atom;
	exports.AWESOME_FACE = awesomeFace;
	exports.BANANA = banana;
	exports.BEARFACE = bearface;
	exports.BEAVIS_ZEN = beavis_zen;
	exports.BEES = bees;
	exports.BILL_THE_CAT = billTheCat;
	exports.BIOHAZARD = biohazard;
	exports.BISHOP = bishop;
	exports.BLACK_MESA = blackMesa;
	exports.BONG = bong;
	exports.BOX = box;
	exports.BROKEN_HEART = brokenHeart;
	exports.BUD_FROGS = budFrogs;
	exports.BUNNY = bunny;
	exports.C3PO = C3PO;
	exports.CAKE = cake;
	exports.CAKE_WITH_CANDLES = cakeWithCandles;
	exports.CAT = cat;
	exports.CAT2 = cat2;
	exports.CATFENCE = catfence;
	exports.CHARIZARDVICE = charizardvice;
	exports.CHARLIE = charlie;
	exports.CHEESE = cheese;
	exports.CHESSMEN = chessmen;
	exports.CHITO = chito;
	exports.CLAW_ARM = clawArm;
	exports.CLIPPY = clippy;
	exports.COMPANION_CUBE = companionCube;
	exports.COWER = cower;
	exports.COWFEE = cowfee;
	exports.CTHULHU_MINI = cthulhuMini;
	exports.CUBE = cube;
	exports.DAEMON = daemon;
	exports.DALEK = dalek;
	exports.DALEK_SHOOTING = dalekShooting;
	exports.DEFAULT = DEFAULT_COW;
	exports.DOCKER_WHALE = dockerWhale;
	exports.DOGE = doge;
	exports.DOLPHIN = dolphin;
	exports.DRAGON = dragon;
	exports.DRAGON_AND_COW = dragonAndCow;
	exports.EBI_FURAI = ebi_furai;
	exports.ELEPHANT = elephant;
	exports.ELEPHANT2 = elephant2;
	exports.ELEPHANT_IN_SNAKE = elephantInSnake;
	exports.EXPLOSION = explosion;
	exports.EYES = eyes;
	exports.FAT_BANANA = fatBanana;
	exports.FAT_COW = fatCow;
	exports.FENCE = fence;
	exports.FIRE = fire;
	exports.FLAMING_SHEEP = flamingSheep;
	exports.FOX = fox;
	exports.GHOST = ghost;
	exports.GHOSTBUSTERS = ghostbusters;
	exports.GLADOS = glados;
	exports.GOAT = goat;
	exports.GOAT2 = goat2;
	exports.GOLDEN_EAGLE = goldenEagle;
	exports.HAND = hand;
	exports.HAPPY_WHALE = happyWhale;
	exports.HEDGEHOG = hedgehog;
	exports.HELLOKITTY = hellokitty;
	exports.HIPPIE = hippie;
	exports.HIYA = hiya;
	exports.HIYOKO = hiyoko;
	exports.HOMER = homer;
	exports.HYPNO = hypno;
	exports.IBM = ibm;
	exports.IWASHI = iwashi;
	exports.JELLYFISH = jellyfish;
	exports.KARL_MARX = karl_marx;
	exports.KILROY = kilroy;
	exports.KING = king;
	exports.KISS = kiss;
	exports.KITTEN = kitten;
	exports.KITTY = kitty;
	exports.KNIGHT = knight;
	exports.KOALA = koala;
	exports.KOSH = kosh;
	exports.LAMB = lamb;
	exports.LAMB2 = lamb2;
	exports.LIGHTBULB = lightbulb;
	exports.LOBSTER = lobster;
	exports.LOLLERSKATES = lollerskates;
	exports.LUKE_KOALA = lukeKoala;
	exports.MAILCHIMP = mailchimp;
	exports.MAZE_RUNNER = mazeRunner;
	exports.MECH_AND_COW = mechAndCow;
	exports.MEOW = meow;
	exports.MILK = milk;
	exports.MINOTAUR = minotaur;
	exports.MONA_LISA = monaLisa;
	exports.MOOFASA = moofasa;
	exports.MOOGHIDJIRAH = mooghidjirah;
	exports.MOOJIRA = moojira;
	exports.MOOSE = moose;
	exports.MULE = mule;
	exports.MUTILATED = mutilated;
	exports.NYAN = nyan;
	exports.OCTOPUS = octopus;
	exports.OKAZU = okazu;
	exports.OWL = owl;
	exports.PAWN = pawn;
	exports.PERIODIC_TABLE = periodicTable;
	exports.PERSONALITY_SPHERE = personalitySphere;
	exports.PINBALL_MACHINE = pinballMachine;
	exports.PSYCHIATRICHELP = psychiatrichelp;
	exports.PSYCHIATRICHELP2 = psychiatrichelp2;
	exports.PTERODACTYL = pterodactyl;
	exports.QUEEN = queen;
	exports.R2_D2 = R2D2;
	exports.RADIO = radio;
	exports.REN = ren;
	exports.RENGE = renge;
	exports.ROBOT = robot;
	exports.ROBOTFINDSKITTEN = robotfindskitten;
	exports.ROFLCOPTER = roflcopter;
	exports.ROOK = rook;
	exports.SACHIKO = sachiko;
	exports.SATANIC = satanic;
	exports.SEAHORSE = seahorse;
	exports.SEAHORSE_BIG = seahorseBig;
	exports.SHEEP = sheep;
	exports.SHIKATO = shikato;
	exports.SHRUG = shrug;
	exports.SKELETON = skeleton;
	exports.SMALL = small;
	exports.SMILING_OCTOPUS = smilingOctopus;
	exports.SNOOPY = snoopy;
	exports.SNOOPYHOUSE = snoopyhouse;
	exports.SNOOPYSLEEP = snoopysleep;
	exports.SPIDERCOW = spidercow;
	exports.SQUID = squid;
	exports.SQUIRREL = squirrel;
	exports.STEGOSAURUS = stegosaurus;
	exports.STIMPY = stimpy;
	exports.SUDOWOODO = sudowoodo;
	exports.SUPERMILKER = supermilker;
	exports.SURGERY = surgery;
	exports.TABLEFLIP = tableflip;
	exports.TAXI = taxi;
	exports.TELEBEARS = telebears;
	exports.TEMPLATE = template;
	exports.THREADER = threader;
	exports.THREECUBES = threecubes;
	exports.TOASTER = toaster;
	exports.TORTOISE = tortoise;
	exports.TURKEY = turkey;
	exports.TURTLE = turtle;
	exports.TUX = tux;
	exports.TUX_BIG = tuxBig;
	exports.TWEETY_BIRD = tweetyBird;
	exports.USA = USA;
	exports.VADER = vader;
	exports.VADER_KOALA = vaderKoala;
	exports.WEEPING_ANGEL = weepingAngel;
	exports.WHALE = whale;
	exports.WIZARD = wizard;
	exports.WOOD = wood;
	exports.WORLD = world;
	exports.WWW = www;
	exports.YASUNA_01 = yasuna_01;
	exports.YASUNA_02 = yasuna_02;
	exports.YASUNA_03 = yasuna_03;
	exports.YASUNA_03A = yasuna_03a;
	exports.YASUNA_04 = yasuna_04;
	exports.YASUNA_05 = yasuna_05;
	exports.YASUNA_06 = yasuna_06;
	exports.YASUNA_07 = yasuna_07;
	exports.YASUNA_08 = yasuna_08;
	exports.YASUNA_09 = yasuna_09;
	exports.YASUNA_10 = yasuna_10;
	exports.YASUNA_11 = yasuna_11;
	exports.YASUNA_12 = yasuna_12;
	exports.YASUNA_13 = yasuna_13;
	exports.YASUNA_14 = yasuna_14;
	exports.YASUNA_16 = yasuna_16;
	exports.YASUNA_17 = yasuna_17;
	exports.YASUNA_18 = yasuna_18;
	exports.YASUNA_19 = yasuna_19;
	exports.YASUNA_20 = yasuna_20;
	exports.YMD_UDON = ymd_udon;
	exports.ZEN_NOH_MILK = zenNohMilk;
	exports.say = say$1;
	exports.think = think$1;

	Object.defineProperty(exports, '__esModule', { value: true });

})));

},{}]},{},[1]);

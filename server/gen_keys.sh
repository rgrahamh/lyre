# Where the keys should be generated
KEY_DIR=".lyre"

# Key details
COUNTRY="US"
STATE="Maryland"
CITY="Gaithersburg"
COMPANY="Nobody"
SECTION="Games"
HOSTNAME="*"
SUBJECT_ALT_NAME="DNS:example.com, IP:[::1], IP:0.0.0.0, DNS:localhost"

if [[ ! -d $KEY_DIR ]]; then
    mkdir -p $KEY_DIR
fi

which openssl &>/dev/null
if [[ $? != 0 ]]; then
    which yum &>/dev/null
    # TODO: Look up package in other distros
    if [[ $? == 0 ]]; then
        yum install openssl openssl-devel
    fi
fi

if [[ ! -f "${KEY_DIR}/lyre.key" ]] && [[ ! -f "${KEY_DIR}/lyre.crt" ]]; then
    openssl req -x509 -newkey rsa:4096 -keyout ${KEY_DIR}/lyre.key -out ${KEY_DIR}/lyre.crt -sha256 -nodes -subj "/C=${COUNTRY}/ST=${STATE}/L=${CITY}/O=${COMPANY}/OU=${SECTION}/CN=${HOSTNAME}" -addext "subjectAltName = ${SUBJECT_ALT_NAME}"
else
    echo "Refusing to overwrite existing certs! If you really want to re-generate them, remove the old ones first."
fi

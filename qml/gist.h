#ifndef GIST_H
#define GIST_H

#include <QWidget>

namespace Ui {
class Gist;
}

class Gist : public QWidget
{
    Q_OBJECT

public:
    explicit Gist(QWidget *parent = 0);
    ~Gist();

private:
    Ui::Gist *ui;
};

#endif // GIST_H
